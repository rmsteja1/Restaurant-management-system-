package routes

import (
	controller "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

 func InvoiceRoutes(incommingRoutes *gin.Engine){
	incommingRoutes.GET("/invoices",controller.GetInvoices())
	incommingRoutes.GET("/invoice/:invoice_id",controller.GetInvoice())
	incommingRoutes.POST("/invoices",controller.CreateInvoice())
	incommingRoutes.PATCH("/invoices:invoices_id",controller.UpdateInvoice())

 }