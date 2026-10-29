package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DirectorSummary struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name string             `bson:"name,omitempty" json:"name,omitempty"`
}

type ActorSummary struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name string             `bson:"name,omitempty" json:"name,omitempty"`
}

type FilmPopulated struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title    string             `bson:"title,omitempty" json:"title,omitempty"`
	Year     int                `bson:"yea,omitemptyr" json:"year,omitempty"`
	Director DirectorSummary    `bson:"director,omitempty" json:"director,omitempty"`
	Actors   []ActorSummary     `bson:"actors,omitempty" json:"actors,omitempty"`
}
