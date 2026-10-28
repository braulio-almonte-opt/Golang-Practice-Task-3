package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Director struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}
