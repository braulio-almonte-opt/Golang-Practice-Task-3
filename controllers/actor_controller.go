package controllers

import (
	"log"
	"net/http"

	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
	"github.com/braulio-almonte-opt/task3-moviesAPI/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActorResponse struct {
	ID   primitive.ObjectID `json:"_id,omitempty"`
	Name string             `json:"name"`
}

type ActorController struct {
	service *service.ActorService
}

func (ctrl *ActorController) RegisterRoutes(router *gin.RouterGroup) {
	actors := router.Group("/actors")
	{
		actors.POST("", ctrl.CreateActor)
		actors.GET("", ctrl.GetAllActors)
		actors.GET("/:id", ctrl.GetActorByID)
		actors.PUT("/:id", ctrl.UpdateActor)
		actors.DELETE("/:id", ctrl.DeleteActor)
	}
}

func NewActorController(service *service.ActorService) *ActorController {
	return &ActorController{
		service: service,
	}
}

func newActorResponse(actor models.Actor) ActorResponse {
	return ActorResponse{
		ID:   actor.ID,
		Name: actor.Name,
	}
}

func (ctrl *ActorController) CreateActor(c *gin.Context) {
	var actor models.Actor
	if err := c.ShouldBindJSON(&actor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdActor, err := ctrl.service.CreateActor(c.Request.Context(), actor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert actor"})
		log.Printf("Error creating actor: %v", err)
		return
	}

	c.JSON(http.StatusCreated, newActorResponse(createdActor))
}

func (ctrl *ActorController) GetActorByID(c *gin.Context) {
	id := c.Param("id")
	actor, err := ctrl.service.GetActorByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Actor not found"})
		return
	}

	c.JSON(http.StatusOK, newActorResponse(actor))
}

func (ctrl *ActorController) GetAllActors(c *gin.Context) {
	actors, err := ctrl.service.ListAllActors(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch actors"})
		return
	}

	response := make([]ActorResponse, len(actors))
	for i, actor := range actors {
		response[i] = newActorResponse(actor)
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *ActorController) UpdateActor(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.service.UpdateActor(c.Request.Context(), id, input.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update actor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Actor updated successfully"})
}

func (ctrl *ActorController) DeleteActor(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.service.DeleteActor(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete actor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Actor deleted successfully"})
}
