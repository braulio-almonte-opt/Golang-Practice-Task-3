package models

// models/film_pop.go (or inside your models package)
import "go.mongodb.org/mongo-driver/bson/primitive"

type DirectorSummary struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type ActorSummary struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type FilmPopulated struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Title    string             `bson:"title" json:"title"`
	Year     int                `bson:"year" json:"year"`
	Director DirectorSummary    `bson:"director" json:"director"`
	Actors   []ActorSummary     `bson:"actors" json:"actors"`
}
