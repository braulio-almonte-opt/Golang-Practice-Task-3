package tests

import (
	"context"
	"testing"

	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
	"github.com/braulio-almonte-opt/task3-moviesAPI/service"
)

func TestActorService_CreateActor_TrimsWhitespace(t *testing.T) {
	actorDAO := &stubActorDAO{actors: map[string]models.Actor{}}
	filmDAO := &stubFilmDAO{}
	actorService := service.NewActorService(actorDAO, filmDAO)

	created, err := actorService.CreateActor(context.Background(), models.Actor{Name: "  Keanu Reeves  "})
	if err != nil {
		t.Fatalf("CreateActor returned an unexpected error: %v", err)
	}
	if created.Name != "Keanu Reeves" {
		t.Fatalf("expected trimmed actor name, got %q", created.Name)
	}
}

func TestActorService_CreateActor_RejectsEmptyName(t *testing.T) {
	actorService := service.NewActorService(&stubActorDAO{actors: map[string]models.Actor{}}, &stubFilmDAO{})

	_, err := actorService.CreateActor(context.Background(), models.Actor{Name: "   "})
	if err == nil {
		t.Fatal("expected error for empty actor name")
	}
}

func TestActorService_GetActorByID_RejectsInvalidID(t *testing.T) {
	actorService := service.NewActorService(&stubActorDAO{actors: map[string]models.Actor{}}, &stubFilmDAO{})

	_, err := actorService.GetActorByID(context.Background(), "nope")
	if err == nil {
		t.Fatal("expected invalid actor ID error")
	}
}
