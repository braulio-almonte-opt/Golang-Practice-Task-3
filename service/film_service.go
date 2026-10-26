package service

import (
	"context"
	"fmt"

	"github.com/braulio-almonte-opt/task3-moviesAPI/dao"
	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
)

type FilmService struct {
	filmDAO     dao.FilmDAO
	directorDAO dao.DirectorDAO
}

func (s *FilmService) GetFilmByDirector(ctx context.Context, directorName string) ([]models.Film, error) {
	director, err := s.directorDAO.GetDirectorByName(ctx, directorName)
	if err != nil {
		return nil, fmt.Errorf("Director not found: %w", err)
	}
	return s.filmDAO.GetFilmByDirector(ctx, director.ID.Hex())
}
