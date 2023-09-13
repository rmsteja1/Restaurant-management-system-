package controllers

import (
	"context"
	"fmt"
	"golang-restaurant-management/database"
	"golang-restaurant-management/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()
var foodCollection *mongo.Collection =database.OpenCollection(database.Client,"food")
var menuCollection *mongo.Collection =database.OpenCollection(database.Client,"menu")


func GetFoods() gin.HandlerFunc{
	return func(c *gin.Context){

	}
}

func GetFood() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel= context.WithTimeout(context.Background(), 100*time.Second)
		foodId := c.Param("food_id")
		var food model.Food
		err:=foodCollection.FindOne(ctx,bson.M{"food_id":foodId}).Decode(&food)
		defer cancel()
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while fetching the food item."})
		}
		c.JSON(http.StatusOK,food)
	}
}

func CreateFood() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		var food model.Food
		var menu model.Menu
		if err:= c.Bind(&food);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		validationError:=validate.Struct(food)
		if validationError!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationError.Error()})
			return
		}
		err:=menuCollection.FindOne(ctx,bson.M{"menu_id":food.Menu_id}).Decode(&menu)
		if err!=nil{
			message:=fmt.Sprintf("menu was not found")
			c.JSON(http.StatusInternalServerError,gin.H{"error":message})
			return
		}
	
		food.Created_at,_ = time.Parse(time.RFC3339, time.Now()).Format(time.RFC3339)
		food.Updated_at,_ = time.Parse(time.RFC3339, time.Now()).Format(time.RFC3339)
		food.ID =primitive.NewObjectID()
		food.Food_id=food.ID.Hex()
		var num = toFixed(*food.Price,2)
		food.Price=&num

		result,insertError:=foodCollection.InsertOne(ctx,food)
		if insertError!=nil{
			message:=fmt.Sprintf("Food item is not inserted")
			c.JSON(http.StatusInternalServerError,gin.H{"error":message})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK,result)
	}
}

func UpdateFood() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func round(num float64) int {

}

func toFixed(num float64,precision int) float64{

}