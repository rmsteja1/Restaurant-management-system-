package controllers

import (
	"context"
	"fmt"
	"golang-restaurant-management/database"
	"golang-restaurant-management/model"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	
		food.Created_at,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_at,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
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
	return func(c *gin.Context) {
		var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var newFood model.Food
		
		if err:=c.Bind(&newFood);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		var updateObj primitive.D

		if newFood.Name!=nil{
			updateObj=append(updateObj, bson.E{"name",newFood.Name})
		}
		if newFood.Price!=nil{
			updateObj=append(updateObj, bson.E{"price",newFood.Price})
		}
		if newFood.Food_image!=nil{
			updateObj=append(updateObj, bson.E{"food_image",newFood.Food_image})
		}
		if newFood.Menu_id!=nil{
			fetchError:=menuCollection.FindOne(ctx,bson.M{"menu_id":newFood.Menu_id})
			if fetchError!=nil{
				msg:=fmt.Sprintf("Menu is not found")
				c.JSON(http.StatusBadRequest,gin.H{"error":msg})
			}
		}
		newFood.Updated_at,_ =time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		newFood.Created_at,_ =time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		updateObj= append(updateObj, bson.E{"updated_at",newFood.Updated_at})
		updateObj= append(updateObj, bson.E{"created_at",newFood.Created_at})
		filter:=bson.M{"food_id":newFood.Food_id}
		upsert:=true
		opt :=options.UpdateOptions{
			Upsert: &upsert,
		}

		result,err:=foodCollection.UpdateOne(ctx,filter,bson.D{{"$set",updateObj}},&opt,)
		if err!=nil{
			errorMessage:="food is not updated"
			c.JSON(http.StatusInternalServerError,gin.H{"error":errorMessage})
			return
		}
		c.JSON(http.StatusOK,result)
	}
}

func round(num float64) int {
return int(num+math.Copysign(0.5,num))
}

func toFixed(num float64,precision int) float64{
	output:=math.Pow(10,float64(precision))
	return float64(round(num*output))/output
}