package tests

import (
	"context"
	"testing"

	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
	"github.com/braulio-almonte-opt/task3-moviesAPI/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDirectorService_CreateDirector_TrimsWhitespace(t *testing.T) {
	directorDAO := &stubDirectorDAO{directors: map[string]models.Director{}}
	filmDAO := &stubFilmDAO{}
	directorService := service.NewDirectorService(directorDAO, filmDAO)

	created, err := directorService.CreateDirector(context.Background(), models.Director{Name: "  Christopher Nolan  "})
	if err != nil {
		t.Fatalf("CreateDirector returned an unexpected error: %v", err)
	}
	if created.Name != "Christopher Nolan" {
		t.Fatalf("expected trimmed director name, got %q", created.Name)
	}
}

func TestDirectorService_UpdateDirector_RejectsEmptyName(t *testing.T) {
	directorService := service.NewDirectorService(&stubDirectorDAO{directors: map[string]models.Director{}}, &stubFilmDAO{})

	err := directorService.UpdateDirector(context.Background(), primitive.NewObjectID().Hex(), "   ")
	if err == nil {
		t.Fatal("expected error when director name is empty")
	}
}
