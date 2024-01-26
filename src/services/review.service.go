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

func GetAverageRating(bookID string) (float64, error) {

	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return 0, err
	}

	pipeline := bson.A{bson.D{{"$match", bson.D{{"_id", objectID}}}}, bson.D{{"$unwind", "$reviews"}}, bson.D{{"$group", bson.D{{"_id", "$_id"}, {"Average", bson.D{{"$avg", "$reviews.rating"}}}}}}}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cursor, err := config.DB.Collection("books").Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	var result struct {
		Average float64 `bson:"Average"`
	}

	if cursor.Next(ctx) {
		err := cursor.Decode(&result)
		if err != nil {
			return 0, err
		}
	}

	return result.Average, nil

}
