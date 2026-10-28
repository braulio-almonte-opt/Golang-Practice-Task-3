package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/braulio-almonte-opt/task3-moviesAPI/dao"
	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DirectorService struct {
	directorDAO dao.DirectorDAO
}

func NewDirectorService(directorDAO dao.DirectorDAO) *DirectorService {
	return &DirectorService{
		directorDAO: directorDAO,
	}
}

func (s *DirectorService) CreateDirector(ctx context.Context, director models.Director) (models.Director, error) {
	trimmedNamne := strings.TrimSpace(director.Name)
	if trimmedNamne == "" {
		return models.Director{}, fmt.Errorf("Empty director name")
	}
	director.Name = trimmedNamne

	return s.directorDAO.InsertDirector(ctx, director)
}

func (s *DirectorService) GetDirectorByID(ctx context.Context, directorID string) (models.Director, error) {
	_, err := primitive.ObjectIDFromHex(directorID)
	if err != nil {
		return models.Director{}, fmt.Errorf("Invalid director ID: %v", directorID)
	}
	return s.directorDAO.GetDirectorByID(ctx, directorID)
}

func (s *DirectorService) ListAllDirectors(ctx context.Context) ([]models.Director, error) {
	return s.directorDAO.ListAll(ctx)
}

func (s *DirectorService) UpdateDirector(ctx context.Context, directorID, directorName string) error {
	_, err := primitive.ObjectIDFromHex(directorID)
	if err != nil {
		return fmt.Errorf("Invalid director ID: %v", directorID)
	}

	directorName = strings.TrimSpace(directorName)
	if directorName == "" {
		return fmt.Errorf("Empty director name")
	}

	return s.directorDAO.UpdateDirector(ctx, directorID, directorName)
}

func (s *DirectorService) DeleteDirector(ctx context.Context, DirectorID string) error {
	_, err := primitive.ObjectIDFromHex(DirectorID)
	if err != nil {
		return fmt.Errorf("Invalid directo ID: %v", DirectorID)
	}
	return s.directorDAO.DeleteDirector(ctx, DirectorID)
}
