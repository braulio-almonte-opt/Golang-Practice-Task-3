package service

import (
	"context"

	"github.com/braulio-almonte-opt/task3-moviesAPI/dao"
	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
)

type DirectorService struct {
	directorDAO dao.DirectorDAO
}

func NewDirectorService(directorDAO dao.DirectorDAO) *DirectorService {
	return &DirectorService{
		directorDAO: directorDAO,
	}
}

func (s *DirectorService) CreateDirector(ctx context.Context, Director models.Director) error {
	return s.directorDAO.InsertDirector(ctx, Director)
}

func (s *DirectorService) GetDirectorByID(ctx context.Context, DirectorID string) (models.Director, error) {
	return s.directorDAO.GetDirectorByID(ctx, DirectorID)
}

func (s *DirectorService) ListAllDirectors(ctx context.Context) ([]models.Director, error) {
	return s.directorDAO.ListAll(ctx)
}

func (s *DirectorService) UpdateDirector(ctx context.Context, DirectorID, DirectorName string) error {
	return s.directorDAO.UpdateDirector(ctx, DirectorID, DirectorName)
}

func (s *DirectorService) DeleteDirector(ctx context.Context, DirectorID string) error {
	return s.directorDAO.DeleteDirector(ctx, DirectorID)
}
