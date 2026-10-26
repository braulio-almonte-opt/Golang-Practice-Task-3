package service

import (
	"context"

	"github.com/braulio-almonte-opt/task3-moviesAPI/dao"
	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
)

type ActorService struct {
	actorDAO dao.ActorDAO
}

func NewActorService(actorDAO dao.ActorDAO) *ActorService {
	return &ActorService{
		actorDAO: actorDAO,
	}
}

func (s *ActorService) CreateActor(ctx context.Context, actor models.Actor) error {
	return s.actorDAO.InsertActor(ctx, actor)
}

func (s *ActorService) GetActorByID(ctx context.Context, actorID string) (models.Actor, error) {
	return s.GetActorByID(ctx, actorID)
}

func (s *ActorService) ListAllActors(ctx context.Context) ([]models.Actor, error) {
	return s.actorDAO.ListAll(ctx)
}

func (s *ActorService) UpdateActor(ctx context.Context, actorID, actorName string) error {
	return s.actorDAO.UpdateActor(ctx, actorID, actorName)
}

func (s *ActorService) DeleteActor(ctx context.Context, actorID string) error {
	return s.actorDAO.DeleteActor(ctx, actorID)
}
