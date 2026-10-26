package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Actor struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name"`
}
