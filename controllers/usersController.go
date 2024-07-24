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

func Signup(ctx *gin.Context) {
	publisher := publishers.NewPublisher(initializers.RmqChannel) 
	// get email/pass/username from body
	var body struct {
		Email string
		Password string
		Username string
	}

	if ctx.Bind(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})

		return
	}

	// hash pw
	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})

		return
	}

	// save user
	user := models.User{Email: body.Email, Password: string(hashed), Username: body.Username}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Email or user name already exists",
		})
	}

	err = publisher.UserRegistered(body.Username)
	if err != nil {
		fmt.Printf("error publishing user registered event: %w", err)
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})

		return
	}
	// find user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})

		return
	}

	// compare password with pw hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})

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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token",
		})

		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("auth", tokenstring, 3600*24*30, "", "", true, true)
	err = publisher.UserLoggedIn(user.Username)
	if err != nil {
		fmt.Printf("error publishing user login event: %w", err)
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
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
			return
		}
		err := publisher.UserSignedout(username)
		if err != nil {
		fmt.Printf("error publishing user logout event: %w", err)
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
		ctx.JSON(http.StatusNotFound, gin.H{
		"message": "No user found",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logged in",
		"username": username,
	})
}