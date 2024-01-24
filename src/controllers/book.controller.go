package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weldonkipchirchir/go/bookshelf-api/src/models"
	"github.com/weldonkipchirchir/go/bookshelf-api/src/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const InternalError string = "Internal Server Error"

func AddBook(c *gin.Context) {
	var book models.Book
	// Ensure that the ID is initialized with a new ObjectID
	book.ID = primitive.NewObjectID()

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.CreateBook(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": InternalError})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": book})
}

func UpdateBook(c *gin.Context) {
	var updatedBook models.Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookID := c.Params.ByName("id")

	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	book, err := services.UpdateBook(objectID, &updatedBook)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func GetBook(c *gin.Context) {
	id := c.Params.ByName("id")

	book, err := services.GetBook(id)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": InternalError})
		return
	}

	c.JSON(http.StatusOK, book)

}

func GetBooks(c *gin.Context) {
	books, err := services.GetBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": InternalError})
		return
	}

	c.JSON(http.StatusOK, books)
}

func SearchBooks(c *gin.Context) {
	query := c.Query("query")
	searchedBooks, err := services.SearchBooks(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": InternalError})
		return
	}

	c.JSON(http.StatusOK, searchedBooks)
}

func GetBooksInProgress(c *gin.Context) {
	booksInProgress, err := services.GetBooksInProgress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": InternalError})
		return
	}

	c.JSON(http.StatusOK, booksInProgress)
}
