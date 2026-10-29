package tests

import (
	"context"
	"errors"

	"github.com/braulio-almonte-opt/task3-moviesAPI/dao"
	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type stubActorDAO struct {
	actors map[string]models.Actor
}

func (s *stubActorDAO) InsertActor(ctx context.Context, actor models.Actor) (models.Actor, error) {
	if actor.ID.IsZero() {
		actor.ID = primitive.NewObjectID()
	}
	s.actors[actor.ID.Hex()] = actor
	return actor, nil
}

func (s *stubActorDAO) GetActorByName(ctx context.Context, actorName string) (models.Actor, error) {
	for _, actor := range s.actors {
		if actor.Name == actorName {
			return actor, nil
		}
	}
	return models.Actor{}, errors.New("actor not found")
}

func (s *stubActorDAO) GetActorByID(ctx context.Context, actorID string) (models.Actor, error) {
	for _, actor := range s.actors {
		if actor.ID.Hex() == actorID {
			return actor, nil
		}
	}
	return models.Actor{}, errors.New("actor not found")
}

func (s *stubActorDAO) FindAllActors(ctx context.Context, actorName string) ([]models.Actor, error) {
	var result []models.Actor
	for _, actor := range s.actors {
		if actor.Name == actorName {
			result = append(result, actor)
		}
	}
	return result, nil
}

func (s *stubActorDAO) ListAll(ctx context.Context) ([]models.Actor, error) {
	result := make([]models.Actor, 0, len(s.actors))
	for _, actor := range s.actors {
		result = append(result, actor)
	}
	return result, nil
}

func (s *stubActorDAO) UpdateActor(ctx context.Context, actorID, actorName string) error {
	for id, actor := range s.actors {
		if id == actorID {
			actor.Name = actorName
			s.actors[id] = actor
			return nil
		}
	}
	return errors.New("actor not found")
}

func (s *stubActorDAO) DeleteActor(ctx context.Context, actorID string) error {
	for id := range s.actors {
		if id == actorID {
			delete(s.actors, id)
			return nil
		}
	}
	return errors.New("actor not found")
}

var _ dao.ActorDAO = (*stubActorDAO)(nil)

type stubDirectorDAO struct {
	directors map[string]models.Director
}

func (s *stubDirectorDAO) InsertDirector(ctx context.Context, director models.Director) (models.Director, error) {
	if director.ID.IsZero() {
		director.ID = primitive.NewObjectID()
	}
	s.directors[director.ID.Hex()] = director
	return director, nil
}

func (s *stubDirectorDAO) GetDirectorByName(ctx context.Context, directorName string) (models.Director, error) {
	for _, director := range s.directors {
		if director.Name == directorName {
			return director, nil
		}
	}
	return models.Director{}, errors.New("director not found")
}

func (s *stubDirectorDAO) GetDirectorByID(ctx context.Context, directorID string) (models.Director, error) {
	for _, director := range s.directors {
		if director.ID.Hex() == directorID {
			return director, nil
		}
	}
	return models.Director{}, errors.New("director not found")
}

func (s *stubDirectorDAO) FindAllDirectors(ctx context.Context, directorName string) ([]models.Director, error) {
	var result []models.Director
	for _, director := range s.directors {
		if director.Name == directorName {
			result = append(result, director)
		}
	}
	return result, nil
}

func (s *stubDirectorDAO) ListAll(ctx context.Context) ([]models.Director, error) {
	result := make([]models.Director, 0, len(s.directors))
	for _, director := range s.directors {
		result = append(result, director)
	}
	return result, nil
}

func (s *stubDirectorDAO) UpdateDirector(ctx context.Context, directorID, directorName string) error {
	for id, director := range s.directors {
		if id == directorID {
			director.Name = directorName
			s.directors[id] = director
			return nil
		}
	}
	return errors.New("director not found")
}

func (s *stubDirectorDAO) DeleteDirector(ctx context.Context, directorID string) error {
	for id := range s.directors {
		if id == directorID {
			delete(s.directors, id)
			return nil
		}
	}
	return errors.New("director not found")
}

var _ dao.DirectorDAO = (*stubDirectorDAO)(nil)

type stubFilmDAO struct {
	films           map[string]models.Film
	directorMatches []models.FilmPopulated
	actorMatches    []models.FilmPopulated
	allPopulated    []models.FilmPopulated
}

func (s *stubFilmDAO) InsertFilm(ctx context.Context, film models.Film) (models.Film, error) {
	if film.ID.IsZero() {
		film.ID = primitive.NewObjectID()
	}
	s.films[film.ID.Hex()] = film
	return film, nil
}

func (s *stubFilmDAO) GetFilmByName(ctx context.Context, filmName string) (models.Film, error) {
	for _, film := range s.films {
		if film.Title == filmName {
			return film, nil
		}
	}
	return models.Film{}, errors.New("film not found")
}

func (s *stubFilmDAO) GetFilmByDirector(ctx context.Context, directorID string) ([]models.Film, error) {
	var result []models.Film
	for _, film := range s.films {
		if film.DirectorID.Hex() == directorID {
			result = append(result, film)
		}
	}
	return result, nil
}

func (s *stubFilmDAO) GetFilmByActor(ctx context.Context, actorID string) ([]models.Film, error) {
	var result []models.Film
	for _, film := range s.films {
		for _, currentActorID := range film.ActorIDs {
			if currentActorID.Hex() == actorID {
				result = append(result, film)
				break
			}
		}
	}
	return result, nil
}

func (s *stubFilmDAO) GetFilmByID(ctx context.Context, filmID string) (models.Film, error) {
	for _, film := range s.films {
		if film.ID.Hex() == filmID {
			return film, nil
		}
	}
	return models.Film{}, errors.New("film not found")
}

func (s *stubFilmDAO) ListAll(ctx context.Context) ([]models.Film, error) {
	result := make([]models.Film, 0, len(s.films))
	for _, film := range s.films {
		result = append(result, film)
	}
	return result, nil
}

func (s *stubFilmDAO) UpdateFilm(ctx context.Context, filmID, filmName string, year int, directorIDStr string, actorIDStrs []string) error {
	for id, film := range s.films {
		if id == filmID {
			film.Title = filmName
			film.Year = year
			film.DirectorID, _ = primitive.ObjectIDFromHex(directorIDStr)
			var actorIDs []primitive.ObjectID
			for _, actorIDStr := range actorIDStrs {
				actorID, err := primitive.ObjectIDFromHex(actorIDStr)
				if err != nil {
					return err
				}
				actorIDs = append(actorIDs, actorID)
			}
			film.ActorIDs = actorIDs
			s.films[id] = film
			return nil
		}
	}
	return errors.New("film not found")
}

func (s *stubFilmDAO) DeleteFilm(ctx context.Context, filmID string) error {
	for id := range s.films {
		if id == filmID {
			delete(s.films, id)
			return nil
		}
	}
	return errors.New("film not found")
}

func (s *stubFilmDAO) GetAllFilmsPopulated(ctx context.Context) ([]models.FilmPopulated, error) {
	return s.allPopulated, nil
}

func (s *stubFilmDAO) GetFilmsByDirectorName(ctx context.Context, directorName string, _ models.SearchOptions) ([]models.FilmPopulated, error) {
	return s.directorMatches, nil
}

func (s *stubFilmDAO) GetFilmsByActorName(ctx context.Context, actorName string, _ models.SearchOptions) ([]models.FilmPopulated, error) {
	return s.actorMatches, nil
}

var _ dao.FilmDAO = (*stubFilmDAO)(nil)
