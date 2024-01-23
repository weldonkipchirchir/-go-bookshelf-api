package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Title           string             `json:"title" validate:"required"`
	Author          string             `json:"author" validate:"required,min=3"`
	Genre           string             `json:"genre" validate:"required,min=3"`
	PublicationYear int                `json:"publication_year" validate:"required"`
	Category        string             `json:"category" validate:"required"`
	Progress        float64            `json:"progress"`
	Reviews         []Review           `json:"reviews"`
	// Add other fields as needed
}

type Review struct {
	ID      primitive.ObjectID `gorm:"primary_key" json:"id"`
	BookID  primitive.ObjectID `json:"book_id"`
	UserID  primitive.ObjectID `json:"user_id"`
	Rating  int                `json:"rating"`
	Comment string             `json:"comment"`
}
