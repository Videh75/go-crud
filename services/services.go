package services

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"main.go/db"
	"main.go/models"
)

func InsertBook(c *gin.Context) {
	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect to server",
		})
	}
	coll := db.Client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("DB_COLLECTION"))
	var doc models.Books
	if err := c.BindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := coll.InsertOne(context.TODO(), doc)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func GetAllBooks(c *gin.Context) {
	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect to server",
		})
	}
	coll := db.Client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("DB_COLLECTION"))

	filter := bson.D{{}}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []models.Books

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"books": results,
	})
}

func InsertManyBooks(c *gin.Context) {
	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect to server",
		})
	}
	coll := db.Client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("DB_COLLECTION"))
	var docs []models.Books
	if err := c.BindJSON(&docs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var doc []interface{}

	for _, book := range docs {
		doc = append(doc, book)
	}

	result, err := coll.InsertMany(context.TODO(), doc)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Documents inserted: %v\n", len(result.InsertedIDs))
	for _, id := range result.InsertedIDs {
		fmt.Printf("Inserted document with _id: %v\n", id)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully inserted %d books", len(result.InsertedIDs)),
	})
}

func UpdateBook(c *gin.Context) {
	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database client is not initialized"})
		return
	}
	coll := db.Client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("DB_COLLECTION"))
	author := c.Query("author")
	filter := bson.D{{Key: "author", Value: author}}

	// Fetch the existing document
	var existingDoc models.Books
	err := coll.FindOne(context.TODO(), filter).Decode(&existingDoc)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	var payload models.Books
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if payload.Title != "" {
		existingDoc.Title = payload.Title
	}
	if payload.Author != "" {
		existingDoc.Author = payload.Author
	}
	update := bson.D{{Key: "$set", Value: existingDoc}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Updated 1 record",
	})
}

func DeleteOneBook(c *gin.Context) {
	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect to server",
		})
	}
	coll := db.Client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("DB_COLLECTION"))
	author := c.Query("author")
	filter := bson.D{{Key: "author", Value: author}}

	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Number of documents deleted: %d\n", result.DeletedCount)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted 1 book",
	})
}

func DeleteManyBooks(c *gin.Context) {
	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect to server",
		})
	}
	coll := db.Client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("DB_COLLECTION"))
	author := c.Query("author")
	filter := bson.D{{Key: "author", Value: author}}

	result, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Number of documents deleted: %d\n", result.DeletedCount)

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully deleted %d books", result.DeletedCount),
	})
}
