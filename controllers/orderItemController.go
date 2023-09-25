package controllers

import (
	"context"
	"golang-restaurant-management/database"
	"golang-restaurant-management/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderItemPack struct {
	Table_id string
	Order_items []model.OrderItem
}

var orderItemCollection *mongo.Collection= database.OpenCollection(database.Client,"orderItem")

func GetOrderItems() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel =context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		results,fetchError:=orderItemCollection.Find(context.TODO(),bson.M{})
		if fetchError!=nil{
			erroMessage:="Couldn't fetch order items."
			c.JSON(http.StatusInternalServerError,gin.H{"error":erroMessage})
			return
		}
		var allOrderItems []bson.M
		if err :=results.All(ctx,&allOrderItems); err!=nil{
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusSeeOther,allOrderItems)
	}
}

func GetOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel =context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var orderItem model.OrderItem
		var orderItemId =c.Param("orderItem_id")

		err:=orderItemCollection.FindOne(ctx,bson.M{"Order_item_id":orderItemId}).Decode(&orderItem)
		if err!=nil{
			errMessage:="error in fetching the order item"
			c.JSON(http.StatusInternalServerError,gin.H{"error":errMessage})
			return
		}
		c.JSON(http.StatusOK,orderItem)
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		orderId:=c.Param("order_id")
		allOrderItems,err :=ItemsByOrder(orderId)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error while listening to order items by order."})
			return
		}
		c.JSON(http.StatusOK,allOrderItems)
	}
}

func ItemsByOrder (orderId string) (orderItems []primitive.M,err error){

}

func CreateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		var orderItemPack OrderItemPack
		var order model.Order
		
		if err:=c.Bind(&orderItemPack);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		order.Order_date,_=time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		orderItemsToBeInserted:=[]interface {}{}
		order.Table_id = &orderItemPack.Table_id
		order_id:=orderItemOrderCreator(order)

		for _,orderItem := range orderItemPack.Order_items{
			orderItem.Order_id=order_id
			validationError:=validate.Struct(orderItem)

			if validationError!=nil{
				c.JSON(http.StatusBadRequest,gin.H{"error":validationError})
				return
			}
			orderItem.ID=primitive.NewObjectID()
			orderItem.Created_at,_= time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
			orderItem.Updated_at,_= time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
			orderItem.Order_item_id = orderItem.ID.Hex()
			var num = toFixed(*orderItem.Unit_price,2)
			orderItem.Unit_price=&num
			orderItemsToBeInserted=append(orderItemsToBeInserted, orderItem)
		}

		insertOrderItems,err := orderItemCollection.InsertMany(ctx,orderItemsToBeInserted)
		if err!=nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK,insertOrderItems)	
		
	}
}

func UpdateOrderItem() gin.HandlerFunc{
	return func(c*gin.Context){
		var ctx,cancel =context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var updateObject primitive.D
		var orderItem model.OrderItem
		orderItemId:=c.Param("order_item_id")
		filter:=bson.M{"order_item_id":orderItemId}

		if orderItem.Unit_price!=nil{
			updateObject = append(updateObject, bson.E{"unit_price",*&orderItem.Unit_price})
		}
		if orderItem.Quantity!=nil{
			updateObject=append(updateObject, bson.E{"quantity",*&orderItem.Quantity})
		}
		if orderItem.Food_id!=nil{
			updateObject=append(updateObject, bson.E{"food_id",*&orderItem.Food_id})
		}

		orderItem.Updated_at,_= time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		updateObject =append(updateObject,bson.E{"updated_at",orderItem.Updated_at})

		upsert:=true
		opt:=options.UpdateOptions{
			Upsert:&upsert,
		}
		result,insertErr:=orderItemCollection.UpdateOne(ctx,
		filter,
	bson.D{{"$set",updateObject},},&opt)
	if insertErr!=nil{
		errMessage:="order item updation fail"
		c.JSON(http.StatusInternalServerError,gin.H{"error":errMessage})
		return
	}
	c.JSON(http.StatusOK,result)
	}
}

