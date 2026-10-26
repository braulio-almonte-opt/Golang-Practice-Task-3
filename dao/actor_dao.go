package dao

import (
	"context"
	"fmt"

	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ActorDAO interface {
	InsertActor(ctx context.Context, actor models.Actor) error
	GetActorByName(ctx context.Context, actorName string) (models.Actor, error)
	GetActorByID(ctx context.Context, actorID string) (models.Actor, error)
	FindAllActors(ctx context.Context, actorName string) ([]models.Actor, error)
	ListAll(ctx context.Context) ([]models.Actor, error)
	UpdateActor(ctx context.Context, actorID, actorName string) error
	DeleteActor(ctx context.Context, actorID string) error
}

type MongoActorDAO struct {
	collection *mongo.Collection
}

func NewActorDAO(client *mongo.Client, dbName string) *MongoActorDAO {
	return &MongoActorDAO{
		collection: client.Database(dbName).Collection("actors"),
	}
}

var _ ActorDAO = (*MongoActorDAO)(nil)

func (dao *MongoActorDAO) InsertActor(ctx context.Context, actor models.Actor) error {
	result, err := dao.collection.InsertOne(ctx, actor)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		actor.ID = oid
	}

	return nil
}

func (dao *MongoActorDAO) GetActorByID(ctx context.Context, actorID string) (models.Actor, error) {
	var actor models.Actor
	id, err := primitive.ObjectIDFromHex(actorID)
	if err != nil {
		return actor, err
	}
	filter := bson.M{"_id": id}
	err = dao.collection.FindOne(ctx, filter).Decode(&actor)
	if err != nil {
		return actor, err
	}
	return actor, nil
}

func (dao *MongoActorDAO) GetActorByName(ctx context.Context, actorName string) (models.Actor, error) {
	var actor models.Actor
	filter := bson.D{{Key: "name", Value: actorName}}
	err := dao.collection.FindOne(ctx, filter).Decode(&actor)
	return actor, err
}

func (dao *MongoActorDAO) FindAllActors(ctx context.Context, actorName string) ([]models.Actor, error) {
	var actors []models.Actor
	filter := bson.D{{Key: "name", Value: actorName}}

	cursor, err := dao.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &actors); err != nil {
		return nil, err
	}

	return actors, nil
}

func (dao *MongoActorDAO) ListAll(ctx context.Context) ([]models.Actor, error) {
	var actors []models.Actor
	filter := bson.M{}
	cursor, err := dao.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &actors); err != nil {
		return nil, err
	}

	return actors, nil
}

func (dao *MongoActorDAO) UpdateActor(ctx context.Context, actorID, actorName string) error {
	id, err := primitive.ObjectIDFromHex(actorID)

	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{"name": actorName}}
	filter := bson.M{"_id": id}
	result, err := dao.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("New result:", result)
	return nil
}

func (dao *MongoActorDAO) DeleteActor(ctx context.Context, actorID string) error {
	id, err := primitive.ObjectIDFromHex(actorID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}
	result, err := dao.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	fmt.Println("Deleted actor: ", result)
	return nil
}
