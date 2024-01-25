// services/reviewService.go
package services

import (
	"context"
	"time"

	"github.com/weldonkipchirchir/go/bookshelf-api/src/config"
	"github.com/weldonkipchirchir/go/bookshelf-api/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddReview(bookID string, userID primitive.ObjectID, review *models.Review) (*models.Review, error) {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, err
	}

	review.BookID = objectID
	review.UserID = userID
	review.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": objectID}
	update := bson.M{"$push": bson.M{"reviews": review}}

	result, err := config.DB.Collection("books").UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.ModifiedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return review, nil
}
