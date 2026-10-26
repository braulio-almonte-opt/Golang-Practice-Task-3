package service

import (
	"context"

	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
)

type ActorService struct {
	actors models.Actor
}

func (s *ActorService) GetActorByID(ctx context.Context, actorID string) (*models.Actor, error) {
	return s.GetActorByID(ctx, actorID)
}
