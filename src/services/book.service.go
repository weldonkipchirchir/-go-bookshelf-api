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

func CreateBook(book *models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err := config.DB.Collection("books").InsertOne(ctx, book)
	defer cancel()
	return err
}

func UpdateBook(bookID primitive.ObjectID, updateBook *models.Book) (*models.Book, error) {
	filter := bson.M{"_id": bookID}

	update := bson.M{"$set": updateBook}

	updateBook.ID = bookID

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result, err := config.DB.Collection("books").UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	// check if any document has been modified
	if result.ModifiedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return updateBook, nil
}

func GetBook(id string) (*models.Book, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var book models.Book

	err = config.DB.Collection("books").FindOne(ctx, bson.M{"_id": objectID}).Decode(&book)

	if err != nil {
		return nil, err
	}

	return &book, nil

}

func GetBooks() ([]models.Book, error) {

	var books []models.Book

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cursor, err := config.DB.Collection("books").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &books)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func SearchBooks(query string) ([]models.Book, error) {
	var SearchBooks []models.Book
	filter := bson.M{
		"$or": []bson.M{
			{"title": primitive.Regex{Pattern: query, Options: "i"}},
			{"Author": primitive.Regex{Pattern: query, Options: "i"}},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	cursor, err := config.DB.Collection("books").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	cursor.All(ctx, &SearchBooks)
	if err != nil {
		return nil, err
	}
	return SearchBooks, nil
}

func GetBooksInProgress() ([]models.Book, error) {
	var booksInProgress []models.Book
	filter := bson.M{"progress": bson.M{"$gt": 0}}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cursor, err := config.DB.Collection("books").Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &booksInProgress)
	if err != nil {
		return nil, err
	}

	return booksInProgress, nil
}
