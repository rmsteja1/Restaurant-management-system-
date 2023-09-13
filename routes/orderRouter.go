package routes

import (
	controller "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(inputRoutes *gin.Engine) {
	inputRoutes.GET("/order",controller.GetOrders())
	inputRoutes.GET("/orders/:order_id",controller.GetOrder())
	inputRoutes.POST("/orders",controller.CreateOrder())
	inputRoutes.PATCH("/orders/:order_id",controller.UpdateOrder())
}