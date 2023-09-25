package controllers

import (
	"context"
	"fmt"
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

type invoiceViewFormat struct{
	Invoice_id			string
	Payment_method		string
	Order_id			string
	Payment_status		*string
	Payment_due			interface{}
	Table_number		interface{}
	Payment_due_date	time.Time
	order_details		interface{}
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client,"invoice")

func GetInvoices () gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var allInvoices []bson.M
		result,dataFetchError:=invoiceCollection.Find(context.TODO(),bson.M{})
		if dataFetchError!=nil{
			log.Fatal(dataFetchError)
			return
		}
		if err:=result.All(ctx,&allInvoices);err!=nil{
			c.JSON(http.StatusInternalServerError,err.Error())
			return
		}
		c.JSON(http.StatusOK,result)
	}
}

func GetInvoice () gin.HandlerFunc{
	return func (c *gin.Context){
		var ctx,cancel =context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		invoiceId:=c.Param("invoice_id")
		if invoiceId==""{
			emptyInputMessage:=fmt.Sprintln("Input should not be empty")
			c.JSON(http.StatusBadRequest,gin.H{"error":emptyInputMessage})
			return
		}
		var newInvoice model.Invoice
		fetchError:=invoiceCollection.FindOne(ctx,bson.M{"invoice_id":invoiceId}).Decode(&newInvoice)
		if fetchError!=nil{
			fetchErrorMessage:="Error occured while searching the invoice."
			c.JSON(http.StatusInternalServerError,gin.H{"error":fetchErrorMessage})
			return
		}
		var invoiceView invoiceViewFormat
		allOrderItems,err:=ItemsByOrder(newInvoice.Order_id)
		invoiceView.Order_id =newInvoice.Order_id
		invoiceView.Payment_due_date =newInvoice.Payment_due_date

		invoiceView.Payment_method="null"
		if newInvoice.Payment_method!=nil{
			invoiceView.Payment_method= *newInvoice.Payment_method
		}

		invoiceView.Invoice_id = newInvoice.Invoice_id
		invoiceView.Payment_status = newInvoice.Payment_status
		invoiceView.Payment_due = allOrderItems[0]["payment_due"]
		invoiceView.Table_number = allOrderItems[0]["table_number"]
		invoiceView.order_details = allOrderItems[0]["order_items"]
		c.JSON(http.StatusOK,invoiceView)
	}
}

func CreateInvoice() gin.HandlerFunc{	
	return func(c *gin.Context){
		var ctx,cancel =context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var newInvoice model.Invoice
		var newOrder model.Order
		var orderId =c.Param("order_id")
		inputErrorMess:="Not a valid input"
		if inputError:=c.Bind(&newInvoice);inputError!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":inputErrorMess})
			return
		}
		validationError:=validate.Struct(newOrder)
		if validationError!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":inputErrorMess})
		}
		if orderId!=""{
			orderError:=orderCollection.FindOne(ctx,bson.M{"order_id":orderId})
			if orderError!=nil{
				c.JSON(http.StatusBadRequest,gin.H{"error":"No a valid order id"})
				return
			}
		}
		status:="PENDING"
		if newInvoice.Payment_status==nil{
			newInvoice.Payment_status=&status
		}

		newInvoice.Payment_due_date,_ = time.Parse(time.RFC3339,time.Now().AddDate(0,0,10).Format(time.RFC3339))
		newInvoice.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		newInvoice.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		newInvoice.ID=primitive.NewObjectID()
		newInvoice.Invoice_id=newInvoice.ID.Hex()

		result,insertErr:=invoiceCollection.InsertOne(ctx,newInvoice)
		if insertErr!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Record not inserted"})
			return
		}
		c.JSON(http.StatusOK,result)
	}
}

func UpdateInvoice() gin.HandlerFunc{
	return func (c *gin.Context){
		var ctx,cancel =context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		invoiceId:=c.Param("invoice_id")
		var newInvoice model.Invoice
		if err:=c.Bind(&newInvoice);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		filter:=bson.M{"invoice_id":invoiceId}
		var updateObj primitive.D
		
		if newInvoice.Payment_method!=nil{
			updateObj=append(updateObj, bson.E{"payment_method",newInvoice.Payment_method})
		}
		if newInvoice.Payment_status!=nil{
			updateObj=append(updateObj,bson.E{"payment_status",newInvoice.Payment_status})
		}
		// if newInvoice.Payment_due_date{
		// 	updateObj=append(updateObj, bson.E{"payment_due_date",newInvoice.Payment_due_date})
		// }
		newInvoice.Updated_at,_= time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		updateObj=append(updateObj, bson.E{"updated_at",newInvoice.Updated_at})
		upsert:=true
		opt:=options.UpdateOptions{
			Upsert: &upsert,
		}

		result,insertError:=invoiceCollection.UpdateOne(ctx,filter,bson.D{{"$set",updateObj},},&opt)
		if insertError!=nil{
			insertErrorMessage:=fmt.Sprintf("Invoice is not updated successfully")
			c.JSON(http.StatusInternalServerError,gin.H{"error":insertErrorMessage})
			return
		}
		c.JSON(http.StatusOK,result)
	}
}


