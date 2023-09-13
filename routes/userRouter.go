package routes

import (
	controller "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("users/",controller.GetUsers())
	incommingRoutes.GET("/user/:food_id",controller.GetUser())
	incommingRoutes.POST("/signUp",controller.SignUp())
	incommingRoutes.POST("/logIn",controller.LogIn())
}