// controllers/reviewController.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weldonkipchirchir/go/bookshelf-api/src/models"
	"github.com/weldonkipchirchir/go/bookshelf-api/src/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddReview(c *gin.Context) {
	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	// Convert userID to primitive.ObjectID
	var userIDObj primitive.ObjectID
	switch v := userID.(type) {
	case primitive.ObjectID:
		userIDObj = v
	case string:
		// Attempt to convert the string to ObjectID
		objectID, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			// Conversion failed, handle the error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user_ID information", "details": err.Error()})
			return
		}
		userIDObj = objectID
	default:
		// Type is not string or ObjectID, handle the error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user Id information"})
		return
	}

	bookID := c.Params.ByName("bookID")

	insertedReview, err := services.AddReview(bookID, userIDObj, &review)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error 1": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, insertedReview)
}

func GetAverageRating(c *gin.Context) {
	_, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	bookId := c.Params.ByName("id")
	rating, err := services.GetAverageRating(bookId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"averageRating": rating})
}
