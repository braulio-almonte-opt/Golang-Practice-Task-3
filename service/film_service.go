package service

import (
	"context"
	"fmt"

	"github.com/braulio-almonte-opt/task3-moviesAPI/dao"
	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilmService struct {
	filmDAO     dao.FilmDAO
	directorDAO dao.DirectorDAO
}

func NewFilmService(filmDAO dao.FilmDAO, directorDAO dao.DirectorDAO) *FilmService {
	return &FilmService{
		filmDAO:     filmDAO,
		directorDAO: directorDAO,
	}
}

func (s *FilmService) CreateFilm(ctx context.Context, title string, year int, directorIDStr string, actorIDStrs []string) (models.Film, error) {

	directorID, err := primitive.ObjectIDFromHex(directorIDStr)
	if err != nil {
		return models.Film{}, fmt.Errorf("invalid director ID: %w", err)
	}

	_, err = s.directorDAO.GetDirectorByID(ctx, directorIDStr)
	if err != nil {
		return models.Film{}, fmt.Errorf("director not found: %w", err)
	}

	actorIDs := make([]primitive.ObjectID, len(actorIDStrs))
	for i, aIDStr := range actorIDStrs {
		aID, err := primitive.ObjectIDFromHex(aIDStr)
		if err != nil {
			return models.Film{}, fmt.Errorf("invalid actor ID at index %d: %w", i, err)
		}
		actorIDs[i] = aID
	}

	film := models.Film{
		Title:      title,
		Year:       year,
		DirectorID: directorID,
		ActorIDs:   actorIDs,
	}

	if err := s.filmDAO.InsertFilm(ctx, film); err != nil {
		return models.Film{}, err
	}

	return film, nil
}

func (s *FilmService) GetFilmByID(ctx context.Context, filmID string) (models.Film, error) {
	return s.filmDAO.GetFilmByID(ctx, filmID)
}

func (s *FilmService) ListAllFilms(ctx context.Context) ([]models.Film, error) {
	return s.filmDAO.ListAll(ctx)
}

func (s *FilmService) UpdateFilm(ctx context.Context, filmID string, title string, year int, directorIDStr string, actorIDStrs []string) error {
	return s.filmDAO.UpdateFilm(ctx, filmID, title, year, directorIDStr, actorIDStrs)
}

func (s *FilmService) DeleteFilm(ctx context.Context, filmID string) error {
	return s.filmDAO.DeleteFilm(ctx, filmID)
}

func (s *FilmService) GetFilmsByDirector(ctx context.Context, directorIDStr string) ([]models.Film, error) {
	return s.filmDAO.GetFilmByDirector(ctx, directorIDStr)
}
