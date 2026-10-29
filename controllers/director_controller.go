package controllers

import (
	"net/http"

	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
	"github.com/braulio-almonte-opt/task3-moviesAPI/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DirectorResponse struct {
	ID   primitive.ObjectID `json:"_id,omitempty"`
	Name string             `json:"name"`
}

type DirectorController struct {
	service *service.DirectorService
}

func (ctrl *DirectorController) RegisterRoutes(router *gin.RouterGroup) {
	directors := router.Group("/directors")
	{
		directors.POST("", ctrl.CreateDirector)
		directors.GET("", ctrl.GetAllDirectors)
		directors.GET("/:id", ctrl.GetDirectorByID)
		directors.PUT("/:id", ctrl.UpdateDirector)
		directors.DELETE("/:id", ctrl.DeleteDirector)
	}
}
func NewDirectorController(service *service.DirectorService) *DirectorController {
	return &DirectorController{
		service: service,
	}
}

func NewDirectorResponse(director models.Director) DirectorResponse {
	return DirectorResponse{
		ID:   director.ID,
		Name: director.Name,
	}
}

func (ctrl *DirectorController) CreateDirector(c *gin.Context) {
	var director models.Director
	if err := c.ShouldBindJSON(&director); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdDirector, err := ctrl.service.CreateDirector(c.Request.Context(), director)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert director"})
		return
	}

	c.JSON(http.StatusCreated, NewDirectorResponse(createdDirector))
}

func (ctrl *DirectorController) GetDirectorByID(c *gin.Context) {
	id := c.Param("id")
	director, err := ctrl.service.GetDirectorByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Director not found"})
		return
	}
	c.JSON(http.StatusOK, NewDirectorResponse(director))
}

func (ctrl *DirectorController) GetAllDirectors(c *gin.Context) {
	directors, err := ctrl.service.ListAllDirectors(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch directors"})
		return
	}

	response := make([]DirectorResponse, len(directors))
	for i, director := range directors {
		response[i] = NewDirectorResponse(director)
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *DirectorController) UpdateDirector(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.service.UpdateDirector(c.Request.Context(), id, input.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update director"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Director updated successfully"})
}

func (ctrl *DirectorController) DeleteDirector(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.service.DeleteDirector(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete director"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Director deleted successfully"})
}
