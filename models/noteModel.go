package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID				primitive.ObjectID	`bson:"_id"`
	Title 			string				`json:"title" validate:"min=1,max=100"`
	Content			string				`json:"content" validate:"min=1,max=1000"`
	User_id			string				`json:"user_id"`
}