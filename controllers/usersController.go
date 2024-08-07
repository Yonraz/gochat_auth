package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yonraz/gochat_auth/events/publishers"
	"github.com/yonraz/gochat_auth/initializers"
	"github.com/yonraz/gochat_auth/models"
	"golang.org/x/crypto/bcrypt"
)
type ErrorResponse struct {
    Errors []string `json:"errors"`
}


func Signup(ctx *gin.Context) {
	publisher := publishers.NewPublisher(initializers.RmqChannel) 
	// get email/pass/username from body
	var body struct {
		Email string
		Password string
		Username string
	}

	if ctx.Bind(&body) != nil {
		response := &ErrorResponse{
			Errors: []string{"Please use valid email, username and password"},
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// hash pw
	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		response := &ErrorResponse{
			Errors: []string{"The server encountered an error, please try again"},
		}
		ctx.JSON(http.StatusInternalServerError, response)

		return
	}

	// save user
	user := models.User{Email: body.Email, Password: string(hashed), Username: body.Username}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		response := &ErrorResponse{
			Errors: []string{"Email or user name already exists"},
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err = publisher.UserRegistered(body.Username)
	if err != nil {
		fmt.Printf("error publishing user registered event: %v", err)
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func Signin(ctx *gin.Context) {
	publisher := publishers.NewPublisher(initializers.RmqChannel)
	//get email and pw
	var body struct {
		Email string
		Password string
	}

	if ctx.Bind(&body) != nil {
		response := &ErrorResponse{
			Errors: []string{"Please enter valid email and password values."},
		}
		ctx.JSON(http.StatusBadRequest, response)

		return
	}
	// find user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		response := &ErrorResponse{
			Errors: []string{"Please enter valid email and password values."},
		}
		ctx.JSON(http.StatusBadRequest, response)
		
		return
	}

	// compare password with pw hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		response := &ErrorResponse{
			Errors: []string{"Please enter valid email and password values."},
		}
		ctx.JSON(http.StatusBadRequest, response)
		
		return
	}

	// gen jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"username": user.Username,
	})
	JWT_KEY := os.Getenv("JWT_KEY")
	tokenstring, err := token.SignedString([]byte(JWT_KEY))

	if err != nil {
		response := &ErrorResponse{
			Errors: []string{"Please enter valid email and password values."},
		}
		ctx.JSON(http.StatusBadRequest, response)
		
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("auth", tokenstring, 3600*24*30, "", "", true, true)
	err = publisher.UserLoggedIn(user.Username)
	if err != nil {
		fmt.Printf("error publishing user login event: %v", err)
	}
	// send back
	ctx.JSON(http.StatusOK, gin.H{
		"email": user.Email,
		"username": user.Username,
	})
}

func Signout(ctx *gin.Context) {
	publisher := publishers.NewPublisher(initializers.RmqChannel)
	// Delete the "auth" cookie
	ctx.SetCookie("auth", "", -1, "/", "", true, true)
	
	// Clear the "currentUserToken" from the context
	user, exists := ctx.Get("currentUser")
	if exists {
		username, ok := user.(string)
		if !ok {
			response := &ErrorResponse{
				Errors: []string{"Please enter valid email and password values."},
				}
			ctx.JSON(http.StatusBadRequest, response)
			return
		}
		err := publisher.UserSignedout(username)
		if err != nil {
		fmt.Printf("error publishing user logout event: %v", err)
		}
		ctx.Set("currentUser", nil)
	}
	ctx.Set("currentUserToken", nil)
	
	// Respond to the client
	ctx.JSON(http.StatusOK, gin.H{"message": "Signed out successfully"})
}

func CurrentUser(ctx *gin.Context) {
	username, exists := ctx.Get("currentUser")
	if !exists {
		response := &ErrorResponse{
			Errors: []string{"No user found."},
		}
		ctx.JSON(http.StatusBadRequest, response)
		
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logged in",
		"username": username,
	})
}
