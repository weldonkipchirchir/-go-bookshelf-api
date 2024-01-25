package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/weldonkipchirchir/go/bookshelf-api/src/controllers"
	"github.com/weldonkipchirchir/go/bookshelf-api/src/middleware"
)

func SetUpReviews(router *gin.Engine) {
	reviewGroup := router.Group("/api/v1/reviews")
	reviewGroup.Use(middleware.Authentication())
	// {
	reviewGroup.POST("/:bookID", controllers.AddReview)
	// 	reviewGroup.GET("/average-rating", controllers.GetAverageRating)
	// }
}
