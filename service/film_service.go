package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/braulio-almonte-opt/task3-moviesAPI/dao"
	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilmService struct {
	filmDAO     dao.FilmDAO
	directorDAO dao.DirectorDAO
	actorDAO    dao.ActorDAO
}

var (
	ErrEmptyID      = errors.New("Empty ID field")
	ErrMissingField = errors.New("There are missing required fields")
)

func NewFilmService(filmDAO dao.FilmDAO, directorDAO dao.DirectorDAO, actorDAO dao.ActorDAO) *FilmService {
	return &FilmService{
		filmDAO:     filmDAO,
		directorDAO: directorDAO,
		actorDAO:    actorDAO,
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
		_, err = s.actorDAO.GetActorByID(ctx, aIDStr)
		if err != nil {
			return models.Film{}, fmt.Errorf("No actor under ID: %v", aIDStr)
		}
		actorIDs[i] = aID
	}

	film := models.Film{
		Title:      title,
		Year:       year,
		DirectorID: directorID,
		ActorIDs:   actorIDs,
	}

	createdFilm, err := s.filmDAO.InsertFilm(ctx, film)
	if err != nil {
		return models.Film{}, err
	}

	return createdFilm, nil
}

func (s *FilmService) GetFilmByID(ctx context.Context, filmID string) (models.Film, error) {

	filmID = strings.TrimSpace(filmID)
	if filmID == "" {
		return models.Film{}, ErrEmptyID
	}
	film, err := s.filmDAO.GetFilmByID(ctx, filmID)
	if err != nil {
		return models.Film{}, fmt.Errorf("Film not found with ID %v", filmID)
	}

	return film, nil
}

func (s *FilmService) ListAllFilms(ctx context.Context) ([]models.Film, error) {
	return s.filmDAO.ListAll(ctx)
}

func (s *FilmService) UpdateFilm(ctx context.Context, filmID string, title string, year int, directorIDStr string, actorIDStrs []string) error {
	_, err := primitive.ObjectIDFromHex(filmID)
	if err != nil {
		return fmt.Errorf("Invalid film ID")
	}

	_, err = primitive.ObjectIDFromHex(directorIDStr)
	if err != nil {
		return fmt.Errorf("No director with ID: %v", directorIDStr)
	}

	stringValidationArray := [2]string{title, directorIDStr}

	for _, s := range stringValidationArray {
		if trimmedString := strings.TrimSpace(s); trimmedString == "" {
			return ErrMissingField
		}
	}

	for _, v := range actorIDStrs {
		_, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			return fmt.Errorf("Invalid actor ID: %v", v)
		}
		_, err = s.actorDAO.GetActorByID(ctx, v)
		if err != nil {
			return fmt.Errorf("No actor under ID: %v", v)
		}
	}

	_, err = s.filmDAO.GetFilmByID(ctx, filmID)
	if err != nil {
		return fmt.Errorf("No film with ID %v", filmID)
	}

	return s.filmDAO.UpdateFilm(ctx, filmID, title, year, directorIDStr, actorIDStrs)
}

func (s *FilmService) DeleteFilm(ctx context.Context, filmID string) error {
	_, err := primitive.ObjectIDFromHex(filmID)
	if err != nil {
		return fmt.Errorf("Invalid film ID %v", filmID)
	}
	return s.filmDAO.DeleteFilm(ctx, filmID)
}

func (s *FilmService) GetFilmsByDirector(ctx context.Context, directorIDStr string) ([]models.Film, error) {
	_, err := primitive.ObjectIDFromHex(directorIDStr)
	if err != nil {
		return []models.Film{}, fmt.Errorf("Invalid director ID: %v", directorIDStr)
	}
	return s.filmDAO.GetFilmByDirector(ctx, directorIDStr)
}

func (s *FilmService) GetFilmsByActor(ctx context.Context, actorIDSts string) ([]models.Film, error) {
	_, err := primitive.ObjectIDFromHex(actorIDSts)
	if err != nil {
		return []models.Film{}, fmt.Errorf("Invalid director ID: %v", actorIDSts)
	}
	return s.filmDAO.GetFilmByActor(ctx, actorIDSts)
}

func (s *FilmService) GetAllFilmsPopulated(ctx context.Context) ([]models.FilmPopulated, error) {
	return s.filmDAO.GetAllFilmsPopulated(ctx)
}

func (s *FilmService) GetFilmsByDirectorName(ctx context.Context, directorName string, opts models.SearchOptions) ([]models.FilmPopulated, error) {
	directorName = strings.TrimSpace(directorName)
	if directorName == "" {
		return []models.FilmPopulated{}, fmt.Errorf("Empty director name")
	}

	_, err := s.directorDAO.GetDirectorByName(ctx, directorName)
	if err != nil {
		return []models.FilmPopulated{}, fmt.Errorf("No director found with name: %v", directorName)
	}

	return s.filmDAO.GetFilmsByDirectorName(ctx, directorName, opts)
}

func (s *FilmService) GetFilmsByActorName(ctx context.Context, actorName string, opts models.SearchOptions) ([]models.FilmPopulated, error) {

	actorName = strings.TrimSpace(actorName)
	if actorName == "" {
		return []models.FilmPopulated{}, fmt.Errorf("Empty actor name")
	}
	_, err := s.actorDAO.GetActorByName(ctx, actorName)
	if err != nil {
		return []models.FilmPopulated{}, fmt.Errorf("No actor with name: %v", actorName)
	}

	return s.filmDAO.GetFilmsByActorName(ctx, actorName, opts)
}
