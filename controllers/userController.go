package controllers

import "github.com/gin-gonic/gin"

func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context){

	}
}

func GetUser() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func SignUp () gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func LogIn() gin.HandlerFunc{
	return func(c *gin.Context){}
}

func HashPassword(password string) string{

}

func VerifyPassword(Userpassword string, providedPAssword string) (string,bool){

}

