package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang-restaurant-management/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//var menuCollection *mongo.Collection= database.OpenCollection(database.Client,"menu")

func GetMenus() gin.HandlerFunc{
	return func(c *gin.Context){
		// ctx,cancel:=context.WithTimeout(context.Background(),100*time.Second)
		// result,err:=menuCollection.Find(context.TODO(),bson.M{})
		// defer cancel()
		// if err!=nil{
		// 	c.JSON(http.StatusInternalServerError,"error: error occured while listening to menu items")
		// }
		// var allMenus []bson.M
		// if err=result.All(ctx,&allMenus);err!=nil{
		// 	log.Fatal(err)
		// }
		// c.JSON(http.StatusOK,allMenus)

		ctx,cancel:=context.WithTimeout(context.Background(),100*time.Second)
		result,fetchError:=menuCollection.Find(context.TODO(),bson.M{})
		defer cancel()
		if fetchError!=nil{
			c.JSON(http.StatusInternalServerError,"error:error occured while listening to the menu items")
			return
		}
		var allMealItems []bson.M
		if err:=result.All(ctx,&allMealItems); err!=nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK,allMealItems)
	}
}

func GetMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		menuId:=c.Param("menu_id")
		var newMenu model.Menu
		defer cancel()
		fetchError:= menuCollection.FindOne(ctx,bson.M{"menu_id":menuId}).Decode(&newMenu)
		if fetchError!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":fetchError})
			return
		}
		c.JSON(http.StatusOK,newMenu)
	}
}

func CreateMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		var newMenu model.Menu
		defer cancel()
		if err:=c.Bind(&newMenu); err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Not a valid request"})
		}

		validationError:=validate.Struct(newMenu)
		if validationError!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationError.Error()})
			return
		}
		newMenu.Created_at,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		newMenu.Updated_at,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		newMenu.ID =primitive.NewObjectID()

		result,insertError:=menuCollection.InsertOne(ctx,newMenu)
		if insertError!=nil{
			errorMessage:=fmt.Sprintf("Menu item is not inserted please try again")
			c.JSON(http.StatusInternalServerError,gin.H{"error":errorMessage})
			return
		}
		c.JSON(http.StatusOK,result)
	}
}

func inTimeSpan(start,end,check time.Time) bool {
	return start.After(time.Now()) && end.After(start)
}

func UpdateMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var newMenu model.Menu

		if err:=c.Bind(&newMenu);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		}
		menuId:= c.Param("menu_id")
		filter:= bson.M{"menu_id":menuId}

		var updateObj primitive.D

		if newMenu.Start_Date!=nil && newMenu.End_Date!=nil{
			if !inTimeSpan(*newMenu.Start_Date,*newMenu.End_Date,time.Now()){
				ms:="Kindly recheck the start and end dates"
				c.JSON(http.StatusBadRequest,gin.H{"error":ms})
				return
			}

			updateObj = append(updateObj, bson.E{"start_date",newMenu.Start_Date})
			updateObj = append(updateObj, bson.E{"end_date",newMenu.End_Date})

			if newMenu.Category!=""{
				updateObj=append(updateObj, bson.E{"category",newMenu.Category})
			}
			if newMenu.Name!=""{
				updateObj=append(updateObj, bson.E{"name", newMenu.Name})
			}
			newMenu.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Local().Format(time.RFC3339))
			updateObj=append(updateObj, bson.E{"updated_at",newMenu.Updated_at})

			upsert:=true

			opt:= options.UpdateOptions{
				Upsert:&upsert,
			}
			result,err:=menuCollection.UpdateOne(ctx,filter,bson.D{{"$set",updateObj}},&opt,)
			if err!=nil{
				msg:="Menu update fail please try again"
				c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			}
			c.JSON(http.StatusOK,result)
		}
	}
}
