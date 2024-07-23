package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yonraz/gochat_auth/controllers"
	"github.com/yonraz/gochat_auth/initializers"
	"github.com/yonraz/gochat_auth/middlewares"
)

func init () {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	
	router := gin.Default()

	router.POST("/api/users/signup", controllers.Signup)
	router.POST("/api/users/signin", controllers.Signin)
	router.POST("/api/users/signout", controllers.Signout)
	router.GET("/api/users/currentuser", middlewares.CurrentUser, middlewares.RequireAuth, controllers.CurrentUser)
	router.Run()
}