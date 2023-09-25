package controllers

import (
	"context"
	"fmt"
	"golang-restaurant-management/database"
	"golang-restaurant-management/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderCollection *mongo.Collection= database.OpenCollection(database.Client,"order")
var tableCollection *mongo.Collection = database.OpenCollection(database.Client,"table")

func GetOrders() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		result,fetchErr:=orderCollection.Find(ctx,bson.M{})
		if fetchErr!=nil{
			errorMessage:="Error in fetching orders"
			c.JSON(http.StatusInternalServerError,gin.H{"error":errorMessage})
		}
		var allOrders []bson.M
		if err:=result.All(ctx,&allOrders);err!=nil{
			c.JSON(http.StatusInternalServerError,err.Error())
		}
		c.JSON(http.StatusOK,result)
	}
}

func GetOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx,cancel:= context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		orderId := c.Param("order_id")
		var newOrder model.Order
		fetchError:=orderCollection.FindOne(ctx,bson.M{"order_id":orderId}).Decode(&newOrder)
		if fetchError!=nil{
			fetchError:=fmt.Sprintf("error in finding this order")
			c.JSON(http.StatusInternalServerError,fetchError)
			return
		}
		c.JSON(http.StatusOK,newOrder)
	}
}

func CreateOrder() gin.HandlerFunc{
	return func(c* gin.Context){
		var newOrder model.Order
		var newTable model.Table
		tableId:=c.Param("order_id")
		ctx,cancel:=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		validationErrorMessage:=fmt.Sprintf("Not a valid request")
		if err:=c.Bind(&newOrder);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErrorMessage})
			return
		}
		validationErr:=validate.Struct(newOrder)
		if validationErr!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErrorMessage})
			return
		}
		if tableId!=""{
			tableFethErr:=tableCollection.FindOne(ctx,bson.M{"table_id":tableId}).Decode(&newTable)
			if tableFethErr!=nil{
				tableFethErrMess:=fmt.Sprintln("Provided table id is invalid")
				c.JSON(http.StatusBadRequest,gin.H{"error":tableFethErrMess})
				return
			}
		}
		newOrder.Created_at,_=time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		newOrder.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		newOrder.ID =primitive.NewObjectID()
		newOrder.Order_id=newOrder.ID.Hex()
		result,recordInsertErr:=orderCollection.InsertOne(ctx,newOrder)
		if recordInsertErr!=nil{
			insertErrorMessage:=fmt.Sprintf("Record is not inserted")
			c.JSON(http.StatusInternalServerError,gin.H{"error":insertErrorMessage})
			return
		}
		c.JSON(http.StatusOK,result)
	}
}

func UpdateOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		var newOrder model.Order
		orderId := c.Param("order_id")
		tableId:=c.Param("table_id")
		var updateObj primitive.D
		if err:= c.Bind(&newOrder);err!=nil{
			errMessage:=fmt.Sprintf("Error in reading the input data.")
			c.JSON(http.StatusInternalServerError,gin.H{"error":errMessage})
			return
		}
		var ctx,cancel =context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		if tableId!=""{
			fetchError:=tableCollection.FindOne(ctx,bson.M{"table_id":orderId})
			if fetchError!=nil{
				errMessage:=fmt.Sprintf("Cannot find the table id")
				c.JSON(http.StatusBadRequest,gin.H{"error":errMessage})
				return
			}
		}
		
		
		newOrder.Updated_at, _ =time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		updateObj=append(updateObj, bson.E{"updated_at",newOrder.Updated_at})
		updateObj=append(updateObj, bson.E{"table_id",newOrder.Table_id})

		filter:=bson.M{"order_id":orderId}
		upsert:=true
		opt:=options.UpdateOptions{
			Upsert: &upsert,
		}

		result,updateError:=orderCollection.UpdateOne(ctx,filter,bson.D{{"$set",updateObj},},&opt,)
		if updateError!=nil{
			updateErrorMessage:=fmt.Sprintf("Cannot update the given order")
			c.JSON(http.StatusInternalServerError,gin.H{"error":updateErrorMessage})
			return
		}
		c.JSON(http.StatusOK,result)
	}
}

func orderItemOrderCreator (order model.Order) string{
	var ctx,cancel =context.WithTimeout(context.Background(),100*time.Second)
	order.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
	order.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id=order.ID.Hex()
	orderCollection.InsertOne(ctx,order)
	defer cancel()

	return order.Order_id
}

func isValidTime(orderTime time.Time) bool{
	return !orderTime.After(time.Now())
}