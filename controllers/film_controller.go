package controllers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
	"github.com/braulio-almonte-opt/task3-moviesAPI/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateFilmRequest struct {
	Title      string   `json:"title" binding:"required"`
	Year       int      `json:"year" binding:"required"`
	DirectorID string   `json:"director_id" binding:"required"`
	ActorIDs   []string `json:"actor_ids"`
}

type FilmResponse struct {
	ID         primitive.ObjectID   `json:"id"`
	Title      string               `json:"title"`
	Year       int                  `json:"year"`
	DirectorID primitive.ObjectID   `json:"director_id"`
	ActorIDs   []primitive.ObjectID `json:"actor_ids"`
}

type FilmController struct {
	service *service.FilmService
}

func NewFilmController(service *service.FilmService) *FilmController {
	return &FilmController{
		service: service,
	}
}

func NewFilmResponse(film models.Film) FilmResponse {
	return FilmResponse{
		ID:         film.ID,
		Title:      film.Title,
		Year:       film.Year,
		DirectorID: film.DirectorID,
		ActorIDs:   film.ActorIDs,
	}
}

func (ctrl *FilmController) CreateFilm(c *gin.Context) {
	var req CreateFilmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Error binding film request: %v", err)
		return
	}

	film, err := ctrl.service.CreateFilm(c.Request.Context(), req.Title, req.Year, req.DirectorID, req.ActorIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Error creating film: %v", err)
		return
	}

	c.JSON(http.StatusCreated, NewFilmResponse(film))
}

func (ctrl *FilmController) GetFilmByID(c *gin.Context) {
	id := c.Param("id")
	film, err := ctrl.service.GetFilmByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Film not found"})
		return
	}
	c.JSON(http.StatusOK, NewFilmResponse(film))
}

func (ctrl *FilmController) GetAllFilms(c *gin.Context) {
	films, err := ctrl.service.ListAllFilms(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch films"})
		return
	}

	response := make([]FilmResponse, len(films))
	for i, film := range films {
		response[i] = NewFilmResponse(film)
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *FilmController) UpdateFilm(c *gin.Context) {
	id := c.Param("id")
	var req CreateFilmRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.service.UpdateFilm(c.Request.Context(), id, req.Title, req.Year, req.DirectorID, req.ActorIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update film"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Film updated successfully"})
}

func (ctrl *FilmController) DeleteFilm(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.service.DeleteFilm(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete film"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Film deleted successfully"})
}

func (ctrl *FilmController) GetFilmsByDirector(c *gin.Context) {
	directorId := c.Param("directorId")
	films, err := ctrl.service.GetFilmsByDirector(c.Request.Context(), directorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch films for director"})
		return
	}

	response := make([]FilmResponse, len(films))
	for i, film := range films {
		response[i] = NewFilmResponse(film)
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *FilmController) GetFilmsByActor(c *gin.Context) {
	actorId := c.Param("actorId")
	films, err := ctrl.service.GetFilmsByActor(c.Request.Context(), actorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch films for actor"})
		return
	}

	response := make([]FilmResponse, len(films))
	for i, film := range films {
		response[i] = NewFilmResponse(film)
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *FilmController) GetAllFilmsPopulated(c *gin.Context) {
	films, err := ctrl.service.GetAllFilmsPopulated(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch populated films"})
		return
	}

	c.JSON(http.StatusOK, films)
}

func (ctrl *FilmController) GetFilmsByDirectorName(c *gin.Context) {
	directorName := c.Query("name")
	if directorName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'name' is required"})
		return
	}

	skipStr := c.DefaultQuery("start", c.DefaultQuery("skip", "0"))
	limitStr := c.DefaultQuery("limit", "0")

	skip, _ := strconv.ParseInt(skipStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	var fields []string
	if f := c.Query("fields"); f != "" {
		fields = strings.Split(f, ",")
	}

	var exclude []string
	if e := c.Query("exclude"); e != "" {
		exclude = strings.Split(e, ",")
	}

	opts := models.SearchOptions{
		Skip:    skip,
		Limit:   limit,
		Fields:  fields,
		Exclude: exclude,
	}

	films, err := ctrl.service.GetFilmsByDirectorName(c.Request.Context(), directorName, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch films by director name"})
		return
	}

	c.JSON(http.StatusOK, films)
}

func (ctrl *FilmController) GetFilmsByActorName(c *gin.Context) {
	actorName := c.Query("name")
	if actorName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'name' is required"})
		return
	}

	// Parse pagination parameters (supports both 'start' and 'skip')
	skipStr := c.DefaultQuery("start", c.DefaultQuery("skip", "0"))
	limitStr := c.DefaultQuery("limit", "0")

	skip, _ := strconv.ParseInt(skipStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	// Parse fields projection (e.g., ?fields=title,year)
	var fields []string
	if f := c.Query("fields"); f != "" {
		fields = strings.Split(f, ",")
	}

	// Parse excluded fields projection (e.g., ?exclude=director)
	var exclude []string
	if e := c.Query("exclude"); e != "" {
		exclude = strings.Split(e, ",")
	}

	opts := models.SearchOptions{
		Skip:    skip,
		Limit:   limit,
		Fields:  fields,
		Exclude: exclude,
	}

	films, err := ctrl.service.GetFilmsByActorName(c.Request.Context(), actorName, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch films by actor name"})
		return
	}

	c.JSON(http.StatusOK, films)
}
