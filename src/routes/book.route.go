package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/weldonkipchirchir/go/bookshelf-api/src/controllers"
	"github.com/weldonkipchirchir/go/bookshelf-api/src/middleware"
)

func SetUpBooks(router *gin.Engine) {
	bookGroup := router.Group("/api/v1/books")
	bookGroup.Use(middleware.Authentication())
	{
		bookGroup.POST("/createbook", controllers.AddBook)
		bookGroup.PUT("/:id", controllers.UpdateBook)
		bookGroup.GET("/:id", controllers.GetBook)
		bookGroup.GET("/", controllers.GetBooks)
		bookGroup.GET("/search", controllers.SearchBooks)
		// bookGroup.GET("/sort", controllers.SortBooks)
		bookGroup.GET("/in-progress", controllers.GetBooksInProgress)
	}
}
