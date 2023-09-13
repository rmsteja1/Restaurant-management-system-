package routes

import (
	controller "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/menus",controller.GetMenus())
	incommingRoutes.GET("/menus/menu_id",controller.GetMenu())
	incommingRoutes.POST("menus/",controller.CreateInvoice())
	incommingRoutes.PATCH("/menus/:menu_id",controller.UpdateInvoice())
}