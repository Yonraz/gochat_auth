# Auth Service

This service is responsible for managing user authentication and authorization through JWT, including storing user connection data. It supports CRUD operations through various http endpoints, and publishes serveral Rabbitmq events.

## Table of Contents

- [Features](#features)
- [Technologies](#technologies)
- [Installation](#installation)
- [API Endpoints](#api-endpoints)
  - [Get Users](#get-users)
- [Event Consumers](#event-consumers)
  - [User Registered Handler](#user-registered-handler)
  - [User Logged-in Handler](#user-logged-in-handler)
  - [User Signed-out Handler](#user-signed-out-handler)
- [Models](#models)
  - [User](#user)
- [Services](#services)
  - [Create User](#create-user)
  - [Login User](#login-user)
  - [Signout User](#signout-user)
- [Consumers](#consumers)
  - [User Registered Producer](#user-registered-handler-1)
  - [User Logged-in Producer](#user-logged-in-handler-1)
  - [User Signed-out Producer](#user-signed-out-handler-1)
- [License](#license)

## Features

- Allows for signin, signup and signout operations as well as a validation check (currentUser).
- JWT based authentication
- Stores user information (username, email, hashsed password).
- Produces RabbitMQ messages.
- Uses PostgreSQL as the database.

## Technologies

- Go
- Gin Framework
- GORM
- RabbitMQ
- PostgreSQL

## Installation

_\*This project is part of a super-module (gochat-app), it's intended to run as a microservice on a k8s cluster. these instructions are for local installation._

1. **Clone the repository:**

   ```bash
   git clone https://github.com/yourusername/gochat_auth.git
   cd gochat_auth
   ```

2. **Set up the database:**

   Ensure you have PostgreSQL installed and create a database for the service.

3. **Set up RabbitMQ:**

   Ensure you have RabbitMQ installed and running.

4. **Set up environment variables:**

   Create a `.env` file with the necessary configuration for your PostgreSQL and RabbitMQ instances.

   ```env
   DB_HOST=localhost
   DB_USER=youruser
   DB_PASSWORD=yourpassword
   DB_NAME=yourdb
   DB_PORT=5432

   RABBITMQ_URL=amqp://guest:guest@localhost:5672/
   ```

5. **Run the service:**

   ```bash
   go run main.go
   ```

## API Endpoints

### Authentication Endpoints

#### Signup

**Endpoint:** `/api/auth/signup`

**Method:** `POST`

**Description:** Registers a new user with email, password, and username.

**Request Body:**

    {
      "email": "user@example.com",
      "password": "password123",
      "username": "username"
    }

**Responses:**

- **200 OK:** Successfully registered the user.
- **400 Bad Request:**
  - Failed to read body.
  - Failed to hash password.
  - Email or username already exists.

**Example Response:**

    {
      "status": "success"
    }

#### Signin

**Endpoint:** `/api/auth/signin`

**Method:** `POST`

**Description:** Authenticates a user with email and password and issues a JWT token.

**Request Body:**

    {
      "email": "user@example.com",
      "password": "password123"
    }

**Responses:**

- **200 OK:** Successfully authenticated the user and issued a JWT token.
- **400 Bad Request:**
  - Failed to read body.
  - Invalid email or password.
  - Failed to create token.

**Example Response:**

    {
      "email": "user@example.com",
      "username": "username"
    }

#### Signout

**Endpoint:** `/api/auth/signout`

**Method:** `POST`

**Description:** Signs out the currently authenticated user by clearing the auth cookie.

**Request Headers:**

- `Cookie`: `auth=<JWT token>`

**Responses:**

- **200 OK:** Successfully signed out the user.
- **500 Internal Server Error:** Invalid user data.

**Example Response:**

    {
      "message": "Signed out successfully"
    }

#### Current User

**Endpoint:** `/api/auth/currentuser`

**Method:** `GET`

**Description:** Retrieves the currently authenticated user's username.

**Request Headers:**

- `Cookie`: `auth=<JWT token>`

**Responses:**

- **200 OK:** Successfully retrieved the current user.
- **404 Not Found:** No user found.

**Example Response:**

    {
      "message": "Logged in",
      "username": "username"
    }

### Middleware

- `CurrentUser`: Middleware that extracts the current user from the JWT token.
- `RequireAuth`: Middleware that ensures the user is authenticated.

## Event Producers

The service produces events to RabbitMQ when user registration, login, and logout events occur. This is done to ensure data consistency among all microservices. The relevant services will listen to those events and use them to update their databases.

### User Registered Producer

- **Queue:** `UserRegistrationQueue`
- **Key:** `UserRegisteredKey`
- **Exchange:** `UserEventsExchange`

This producer publishes user registration events. When a new user registers, the Auth Service emits a `UserRegisteredKey` event to the `UserEventsExchange`.

### User Logged-in Producer

- **Queue:** `UserLoginQueue`
- **Key:** `UserLoggedInKey`
- **Exchange:** `UserEventsExchange`

This producer publishes user login events. When a user logs in, the Auth Service emits a `UserLoggedInKey` event to the `UserEventsExchange`.

### User Signed-out Producer

- **Queue:** `UserSignoutQueue`
- **Key:** `UserSignedoutKey`
- **Exchange:** `UserEventsExchange`

This producer publishes user logout events. When a user logs out, the Auth Service emits a `UserSignedoutKey` event to the `UserEventsExchange`.

## Models

### User

```go
type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Username string `gorm:"unique"`
	Password string
}
```
