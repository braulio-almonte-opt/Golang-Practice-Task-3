package tests

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
	"github.com/braulio-almonte-opt/task3-moviesAPI/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFilmService_CreateFilm_Success(t *testing.T) {
	directorID := primitive.NewObjectID()
	actorOneID := primitive.NewObjectID()
	actorTwoID := primitive.NewObjectID()

	directorDAO := &stubDirectorDAO{directors: map[string]models.Director{
		directorID.Hex(): {ID: directorID, Name: "Christopher Nolan"},
	}}
	actorDAO := &stubActorDAO{actors: map[string]models.Actor{
		actorOneID.Hex(): {ID: actorOneID, Name: "Matthew McConaughey"},
		actorTwoID.Hex(): {ID: actorTwoID, Name: "Anne Hathaway"},
	}}
	filmDAO := &stubFilmDAO{films: map[string]models.Film{}}
	filmService := service.NewFilmService(filmDAO, directorDAO, actorDAO)

	created, err := filmService.CreateFilm(context.Background(), "Interstellar", 2014, directorID.Hex(), []string{actorOneID.Hex(), actorTwoID.Hex()})
	if err != nil {
		t.Fatalf("CreateFilm returned an error: %v", err)
	}
	if created.Title != "Interstellar" {
		t.Fatalf("expected title Interstellar, got %q", created.Title)
	}
	if !reflect.DeepEqual(created.ActorIDs, []primitive.ObjectID{actorOneID, actorTwoID}) {
		t.Fatalf("unexpected actor IDs: %#v", created.ActorIDs)
	}
}

func TestFilmService_GetFilmByID_RejectsEmptyID(t *testing.T) {
	filmService := service.NewFilmService(&stubFilmDAO{films: map[string]models.Film{}}, &stubDirectorDAO{directors: map[string]models.Director{}}, &stubActorDAO{actors: map[string]models.Actor{}})

	_, err := filmService.GetFilmByID(context.Background(), "   ")
	if !errors.Is(err, service.ErrEmptyID) {
		t.Fatalf("expected ErrEmptyID, got %v", err)
	}
}

func TestFilmService_GetFilmsByDirectorName_DelegatesToDAO(t *testing.T) {
	directorID := primitive.NewObjectID()
	directorDAO := &stubDirectorDAO{directors: map[string]models.Director{
		directorID.Hex(): {ID: directorID, Name: "Christopher Nolan"},
	}}
	filmDAO := &stubFilmDAO{films: map[string]models.Film{}, directorMatches: []models.FilmPopulated{{Title: "Interstellar", Year: 2014}}}
	filmService := service.NewFilmService(filmDAO, directorDAO, &stubActorDAO{actors: map[string]models.Actor{}})

	films, err := filmService.GetFilmsByDirectorName(context.Background(), " Christopher Nolan ", models.SearchOptions{})
	if err != nil {
		t.Fatalf("GetFilmsByDirectorName returned an error: %v", err)
	}
	if len(films) != 1 || films[0].Title != "Interstellar" {
		t.Fatalf("unexpected films returned: %#v", films)
	}
}

func TestFilmService_GetFilmsByActorName_DelegatesToDAO(t *testing.T) {
	actorID := primitive.NewObjectID()
	actorDAO := &stubActorDAO{actors: map[string]models.Actor{
		actorID.Hex(): {ID: actorID, Name: "Matthew McConaughey"},
	}}
	filmDAO := &stubFilmDAO{films: map[string]models.Film{}, actorMatches: []models.FilmPopulated{{Title: "Interstellar", Year: 2014}}}
	filmService := service.NewFilmService(filmDAO, &stubDirectorDAO{directors: map[string]models.Director{}}, actorDAO)

	films, err := filmService.GetFilmsByActorName(context.Background(), " Matthew McConaughey ", models.SearchOptions{})
	if err != nil {
		t.Fatalf("GetFilmsByActorName returned an error: %v", err)
	}
	if len(films) != 1 || films[0].Title != "Interstellar" {
		t.Fatalf("unexpected films returned: %#v", films)
	}
}
