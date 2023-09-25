package controllers

import (
	"context"
	"fmt"
	"golang-restaurant-management/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//var tableCollection mongo.Collection = database.OpenCollection(database.Client,"table")

func GetTables() gin.HandlerFunc{
	return func (c*gin.Context){
		var ctx,cancel =context.WithTimeout(context.Background(),100*time.Second)
		var allTables []model.Table
		defer cancel()
		result,fetchErr:=tableCollection.Find(context.TODO(),bson.M{})
		if fetchErr!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":fetchErr})
			return
		}
		if err:=result.All(ctx,&allTables);err!=nil{
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK,allTables)
	}
}

func GetTable() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel =context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		tableId:=c.Param("table_id")
		var newTable model.Table
		err:=tableCollection.FindOne(ctx,bson.M{"tableId":tableId}).Decode(&newTable)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error while fetching table"})
			return
		}
		c.JSON(http.StatusOK,newTable)
	}
}

func CreateTable() gin.HandlerFunc{
	return func(c* gin.Context){
		var ctx,cancel =context.WithTimeout(context.Background(),100*time.Second)
		var newTable model.Table
		defer cancel()
		if err:=c.Bind(&newTable);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		
		validationError:=validate.Struct(newTable)
		if validationError!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationError.Error()})
			return
		}
		newTable.Created_at,_ =time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		newTable.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		newTable.ID=primitive.NewObjectID()
		newTable.Table_id=newTable.ID.Hex()

		result,insertErr:=tableCollection.InsertOne(ctx,newTable)

		if insertErr!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Table item was not created."})
			return
		}
		c.JSON(http.StatusOK,result)
	}
}

func UpdateTable() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel =context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var newTable model.Table
		tableId:=c.Param("table_id")

		if err:=c.Bind(&newTable);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		}
		var updateObject primitive.D

		if newTable.Number_of_guests!=nil{
			updateObject=append(updateObject, bson.E{"number_of_guests",newTable.Number_of_guests})
		}
		if newTable.Table_number!=nil{
			updateObject=append(updateObject, bson.E{"table_number",newTable.Table_number})
		}
		newTable.Updated_at,_=time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		upsert:=true
		opt:=options.UpdateOptions{
			Upsert: &upsert,
		}

		filter:=bson.M{"table_id":tableId}
		result,insertErr:=tableCollection.UpdateOne(ctx,filter,bson.D{{"$set",updateObject},},&opt)
		if insertErr!=nil{
			inserErrMess:=fmt.Sprintf("Table item update failed")
			c.JSON(http.StatusInternalServerError,inserErrMess)
			return
		}
		c.JSON(http.StatusOK,result)
	}
}