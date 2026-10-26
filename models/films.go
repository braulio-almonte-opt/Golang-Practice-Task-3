package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Film struct {
	ID         primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Title      string               `json:"title" bson:"title"`
	Year       int                  `json:"year" bson:"year"`
	DirectorID primitive.ObjectID   `json:"director_id" bson:"director_id"`
	ActorIDs   []primitive.ObjectID `json:"actor_ids" bson:"actor_ids"`
}
