package dao

import (
	"context"
	"fmt"

	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FilmDAO interface {
	InsertFilm(ctx context.Context, film models.Film) error
	GetFilmByName(ctx context.Context, filmName string) (models.Film, error)
	GetFilmByDirector(ctx context.Context, directorID string) ([]models.Film, error)
	GetFilmByActor(ctx context.Context, actorID string) ([]models.Film, error)
	GetFilmByID(ctx context.Context, filmID string) (models.Film, error)
	ListAll(ctx context.Context) ([]models.Film, error)
	UpdateFilm(ctx context.Context, filmID, filmName string, year int, directorsName string, actorNames []string) error
	DeleteFilm(ctx context.Context, filmID string) error
	GetAllFilmsPopulated(ctx context.Context) ([]models.FilmPopulated, error)
	GetFilmsByDirectorName(ctx context.Context, directorName string, opt models.SearchOptions) ([]models.FilmPopulated, error)
	GetFilmsByActorName(ctx context.Context, actorName string, opt models.SearchOptions) ([]models.FilmPopulated, error)
}

type MongoFilmDAO struct {
	collection *mongo.Collection
}

func NewFilmDAO(client *mongo.Client, dbName string) *MongoFilmDAO {
	return &MongoFilmDAO{
		collection: client.Database(dbName).Collection("films"),
	}
}

var _ FilmDAO = (*MongoFilmDAO)(nil)

func (dao *MongoFilmDAO) InsertFilm(ctx context.Context, film models.Film) error {
	result, err := dao.collection.InsertOne(ctx, film)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		film.ID = oid
	}

	return nil
}

func (dao *MongoFilmDAO) GetAllFilmsPopulated(ctx context.Context) ([]models.FilmPopulated, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: "directors"},
				primitive.E{Key: "localField", Value: "director_id"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "director"},
			}},
		},
		bson.D{
			primitive.E{Key: "$unwind", Value: "$director"},
		},
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: "actors"},
				primitive.E{Key: "localField", Value: "actor_ids"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "actors"},
			}},
		},
	}

	cursor, err := dao.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return []models.FilmPopulated{}, err
	}
	defer cursor.Close(ctx)

	var films []models.FilmPopulated
	if err := cursor.All(ctx, &films); err != nil {
		return []models.FilmPopulated{}, err
	}
	return films, nil
}

func (dao *MongoFilmDAO) GetFilmsByDirectorName(ctx context.Context, directorName string, opts models.SearchOptions) ([]models.FilmPopulated, error) {
	pipeline := mongo.Pipeline{
		bson.D{primitive.E{Key: "$lookup", Value: bson.D{
			primitive.E{Key: "from", Value: "directors"},
			primitive.E{Key: "localField", Value: "director_id"},
			primitive.E{Key: "foreignField", Value: "_id"},
			primitive.E{Key: "as", Value: "director"},
		}}},
		bson.D{primitive.E{Key: "$unwind", Value: "$director"}},
		bson.D{primitive.E{Key: "$match", Value: bson.D{
			primitive.E{Key: "director.name", Value: primitive.Regex{Pattern: directorName, Options: "i"}},
		}}},
		bson.D{primitive.E{Key: "$lookup", Value: bson.D{
			primitive.E{Key: "from", Value: "actors"},
			primitive.E{Key: "localField", Value: "actor_ids"},
			primitive.E{Key: "foreignField", Value: "_id"},
			primitive.E{Key: "as", Value: "actors"},
		}}},
	}

	if opts.Skip > 0 {
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$skip", Value: opts.Skip}})
	}

	if opts.Limit > 0 {
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$limit", Value: opts.Limit}})
	}

	if len(opts.Fields) > 0 {
		proj := bson.D{primitive.E{Key: "_id", Value: 1}}
		for _, field := range opts.Fields {
			proj = append(proj, primitive.E{Key: field, Value: 1})
		}
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$project", Value: proj}})
	} else if len(opts.Exclude) > 0 {
		proj := bson.D{}
		for _, field := range opts.Exclude {
			proj = append(proj, primitive.E{Key: field, Value: 0})
		}
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$project", Value: proj}})
	}

	cursor, err := dao.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return []models.FilmPopulated{}, err
	}
	defer cursor.Close(ctx)

	var films []models.FilmPopulated
	if err := cursor.All(ctx, &films); err != nil {
		return []models.FilmPopulated{}, err
	}
	return films, nil
}

func (dao *MongoFilmDAO) GetFilmsByActorName(ctx context.Context, actorName string, opts models.SearchOptions) ([]models.FilmPopulated, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: "directors"},
				primitive.E{Key: "localField", Value: "director_id"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "director"},
			}},
		},
		bson.D{
			primitive.E{Key: "$unwind", Value: "$director"},
		},
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: "actors"},
				primitive.E{Key: "localField", Value: "actor_ids"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "actors"},
			}},
		},
		bson.D{
			primitive.E{Key: "$match", Value: bson.D{
				primitive.E{Key: "actors.name", Value: primitive.Regex{Pattern: actorName, Options: "i"}},
			}},
		},
	}

	if opts.Skip > 0 {
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$skip", Value: opts.Skip}})
	}

	if opts.Limit > 0 {
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$limit", Value: opts.Limit}})
	}

	if len(opts.Fields) > 0 {
		proj := bson.D{primitive.E{Key: "_id", Value: 1}}
		for _, field := range opts.Fields {
			proj = append(proj, primitive.E{Key: field, Value: 1})
		}
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$project", Value: proj}})
	} else if len(opts.Exclude) > 0 {
		proj := bson.D{}
		for _, field := range opts.Exclude {
			proj = append(proj, primitive.E{Key: field, Value: 0})
		}
		pipeline = append(pipeline, bson.D{primitive.E{Key: "$project", Value: proj}})
	}

	cursor, err := dao.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return []models.FilmPopulated{}, err
	}
	defer cursor.Close(ctx)

	var films []models.FilmPopulated
	if err := cursor.All(ctx, &films); err != nil {
		return []models.FilmPopulated{}, err
	}
	return films, nil
}

func (dao *MongoFilmDAO) GetFilmByName(ctx context.Context, filmName string) (models.Film, error) {
	var films models.Film
	filter := bson.D{{Key: "name", Value: filmName}}
	err := dao.collection.FindOne(ctx, filter).Decode(&films)
	return films, err
}

func (dao *MongoFilmDAO) GetFilmByDirector(ctx context.Context, directorID string) ([]models.Film, error) {
	id, err := primitive.ObjectIDFromHex(directorID)
	if err != nil {
		return []models.Film{}, err
	}
	var films []models.Film
	filter := bson.D{{Key: "director_id", Value: id}}

	cursor, err := dao.collection.Find(ctx, filter)
	if err != nil {
		return []models.Film{}, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &films); err != nil {
		return nil, err
	}
	return films, nil
}

func (dao *MongoFilmDAO) GetFilmByActor(ctx context.Context, actorID string) ([]models.Film, error) {
	id, err := primitive.ObjectIDFromHex(actorID)
	if err != nil {
		return []models.Film{}, err
	}
	var films []models.Film

	filter := bson.M{"actor_ids": id}

	cursor, err := dao.collection.Find(ctx, filter)
	if err != nil {
		return []models.Film{}, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &films); err != nil {
		return nil, err
	}
	return films, nil
}

func (dao *MongoFilmDAO) GetFilmByID(ctx context.Context, filmID string) (models.Film, error) {
	var film models.Film
	id, err := primitive.ObjectIDFromHex(filmID)
	if err != nil {
		return models.Film{}, err
	}

	filter := bson.M{"_id": id}
	err = dao.collection.FindOne(ctx, filter).Decode(&film)
	if err != nil {
		return models.Film{}, err
	}

	return film, nil
}

func (dao *MongoFilmDAO) ListAll(ctx context.Context) ([]models.Film, error) {
	var films []models.Film
	filter := bson.M{}
	cursor, err := dao.collection.Find(ctx, filter)
	if err != nil {
		return []models.Film{}, err
	}

	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &films); err != nil {
		return []models.Film{}, err
	}

	return films, nil
}

func (dao *MongoFilmDAO) UpdateFilm(ctx context.Context, filmID string, title string, year int, directorIDStr string, actorIDStrs []string) error {

	id, err := primitive.ObjectIDFromHex(filmID)
	if err != nil {
		return err
	}

	directorID, err := primitive.ObjectIDFromHex(directorIDStr)
	if err != nil {
		return err
	}

	actorIDs := make([]primitive.ObjectID, len(actorIDStrs))
	for i, aIDStr := range actorIDStrs {
		aID, err := primitive.ObjectIDFromHex(aIDStr)
		if err != nil {
			return err
		}
		actorIDs[i] = aID
	}

	update := bson.M{
		"$set": bson.M{
			"title":       title,
			"year":        year,
			"director_id": directorID,
			"actor_ids":   actorIDs,
		},
	}

	filter := bson.M{"_id": id}
	result, err := dao.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("Updated Film result:", result)
	return nil
}

func (dao *MongoFilmDAO) DeleteFilm(ctx context.Context, filmID string) error {
	id, err := primitive.ObjectIDFromHex(filmID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": id}
	result, err := dao.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	fmt.Println("Deleted film: ", result)
	return nil
}
