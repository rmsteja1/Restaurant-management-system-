package main

import (
	"golang-restaurant-management/middleware"
	"golang-restaurant-management/routes"
	"os"
	"golang-restaurant-management/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client,"food")

func main(){
port:=os.Getenv("PORT")

if port==""{
	port="8000"
}
	router:=gin.New()
	router.Use(gin.Logger())
	router.Use(middleware.Authentication())

	routes.UserRoutes(router)
	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.InvoiceRoutes(router)
	routes.OrderItemRoutes(router)

	router.Run(":"+port)

}
