package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/db"
	"main.go/initializers"
	"main.go/services"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {

	db.DbConnection()
	defer db.Client.Disconnect(context.TODO())
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
	r.POST("/create", func(c *gin.Context) {
		services.InsertBook(c)
	})
	r.GET("/getAllBooks", func(c *gin.Context) {
		services.GetAllBooks(c)
	})
	r.POST("/createBooks", func(c *gin.Context) {
		services.InsertManyBooks(c)
	})
	r.PUT("/updateBook", func(c *gin.Context) {
		services.UpdateBook(c)
	})
	r.DELETE("/deleteOneBook", func(c *gin.Context) {
		services.DeleteOneBook(c)
	})
	r.DELETE("/deleteManyBooks", func(c *gin.Context) {
		services.DeleteManyBooks(c)
	})
	r.Run()
}
