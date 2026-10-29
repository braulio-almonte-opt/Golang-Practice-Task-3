package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/braulio-almonte-opt/task3-moviesAPI/dao"
	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActorService struct {
	actorDAO dao.ActorDAO
	filmDAO  dao.FilmDAO
}

func NewActorService(actorDAO dao.ActorDAO, filmDAO dao.FilmDAO) *ActorService {
	return &ActorService{
		actorDAO: actorDAO,
		filmDAO:  filmDAO,
	}
}

func (s *ActorService) CreateActor(ctx context.Context, actor models.Actor) (models.Actor, error) {
	trimmedName := strings.TrimSpace(actor.Name)
	if trimmedName == "" {
		return models.Actor{}, fmt.Errorf("Empty actor name")
	}
	actor.Name = trimmedName

	return s.actorDAO.InsertActor(ctx, actor)
}

func (s *ActorService) GetActorByID(ctx context.Context, actorID string) (models.Actor, error) {
	_, err := primitive.ObjectIDFromHex(actorID)
	if err != nil {
		return models.Actor{}, fmt.Errorf("Invalid actor ID: %v", actorID)
	}
	return s.actorDAO.GetActorByID(ctx, actorID)
}

func (s *ActorService) ListAllActors(ctx context.Context) ([]models.Actor, error) {
	return s.actorDAO.ListAll(ctx)
}

func (s *ActorService) UpdateActor(ctx context.Context, actorID, actorName string) error {
	_, err := primitive.ObjectIDFromHex(actorID)
	if err != nil {
		return fmt.Errorf("Invalid actor ID: %v", actorID)
	}

	actorName = strings.TrimSpace(actorName)
	if actorName == "" {
		return fmt.Errorf("Empty actor name")
	}
	return s.actorDAO.UpdateActor(ctx, actorID, actorName)
}

func (s *ActorService) DeleteActor(ctx context.Context, actorID string) error {
	_, err := primitive.ObjectIDFromHex(actorID)
	if err != nil {
		return fmt.Errorf("Invalid actor ID")
	}

	films, err := s.filmDAO.GetFilmByActor(ctx, actorID)
	if err != nil {
		return err
	}

	if len(films) > 0 {
		return fmt.Errorf("Cannot delete actor linked to films: %v", films)
	}

	return s.actorDAO.DeleteActor(ctx, actorID)
}
