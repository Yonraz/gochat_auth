package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yonraz/gochat_auth/controllers"
	"github.com/yonraz/gochat_auth/initializers"
	"github.com/yonraz/gochat_auth/middlewares"
)

func init () {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
	initializers.ConnectToRabbitmq()
}

func main() {
	
	router := gin.Default()

	// Defer closure of the channel and connection
	defer func() {
		if err := initializers.RmqChannel.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ channel: %v", err)
		}
	}()
	defer func() {
		if err := initializers.RmqConn.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %v", err)
		}
	}()

	router.POST("/api/auth/signup", controllers.Signup)
	router.POST("/api/auth/signin", controllers.Signin)
	router.POST("/api/auth/signout", controllers.Signout)
	router.GET("/api/auth/currentuser", middlewares.CurrentUser, middlewares.RequireAuth, controllers.CurrentUser)
	router.Run()
}