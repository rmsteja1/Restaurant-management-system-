package routes

import (
	controller "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func FoodRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/getFoods",controller.GetFoods())
	incommingRoutes.GET("/getFood/:food_id",controller.GetFood())
	incommingRoutes.POST("/foods",controller.CreateFood())
	incommingRoutes.PATCH("/foods/:food_id",controller.UpdateFood())
}