package dao

import (
	"context"
	"fmt"

	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DirectorDAO interface {
	InsertDirector(ctx context.Context, director models.Director) (models.Director, error)
	GetDirectorByName(ctx context.Context, directorName string) (models.Director, error)
	GetDirectorByID(ctx context.Context, directorID string) (models.Director, error)
	FindAllDirectors(ctx context.Context, directorName string) ([]models.Director, error)
	ListAll(ctx context.Context) ([]models.Director, error)
	UpdateDirector(ctx context.Context, directorID, directorName string) error
	DeleteDirector(ctx context.Context, directorID string) error
}

type MongoDirectorDAO struct {
	collection *mongo.Collection
}

func NewDirectorDAO(client *mongo.Client, dbName string) *MongoDirectorDAO {
	return &MongoDirectorDAO{
		collection: client.Database(dbName).Collection("directors"),
	}
}

var _ DirectorDAO = (*MongoDirectorDAO)(nil)

func (dao *MongoDirectorDAO) InsertDirector(ctx context.Context, director models.Director) (models.Director, error) {
	result, err := dao.collection.InsertOne(ctx, director)
	if err != nil {
		return models.Director{}, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		director.ID = oid
	}

	return director, nil
}

func (dao *MongoDirectorDAO) GetDirectorByID(ctx context.Context, directorID string) (models.Director, error) {
	var director models.Director
	id, err := primitive.ObjectIDFromHex(directorID)
	if err != nil {
		return models.Director{}, err
	}
	filter := bson.M{"_id": id}
	err = dao.collection.FindOne(ctx, filter).Decode(&director)
	if err != nil {
		return models.Director{}, err
	}
	return director, nil
}

func (dao *MongoDirectorDAO) GetDirectorByName(ctx context.Context, directorName string) (models.Director, error) {
	var director models.Director
	filter := bson.D{{Key: "name", Value: directorName}}
	err := dao.collection.FindOne(ctx, filter).Decode(&director)
	return director, err
}

func (dao *MongoDirectorDAO) FindAllDirectors(ctx context.Context, directorName string) ([]models.Director, error) {
	var directors []models.Director
	filter := bson.D{{Key: "name", Value: directorName}}

	cursor, err := dao.collection.Find(ctx, filter)
	if err != nil {
		return []models.Director{}, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &directors); err != nil {
		return []models.Director{}, err
	}

	return directors, nil
}

func (dao *MongoDirectorDAO) ListAll(ctx context.Context) ([]models.Director, error) {
	var directors []models.Director
	filter := bson.M{}
	cursor, err := dao.collection.Find(ctx, filter)
	if err != nil {
		return []models.Director{}, err
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &directors); err != nil {
		return []models.Director{}, err
	}

	return directors, nil
}

func (dao *MongoDirectorDAO) UpdateDirector(ctx context.Context, directorID, directorName string) error {
	id, err := primitive.ObjectIDFromHex(directorID)

	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{"name": directorName}}
	filter := bson.M{"_id": id}
	result, err := dao.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("New result:", result)
	return nil
}

func (dao *MongoDirectorDAO) DeleteDirector(ctx context.Context, directorID string) error {
	id, err := primitive.ObjectIDFromHex(directorID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}
	result, err := dao.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	fmt.Println("Deleted Director: ", result)
	return nil
}
