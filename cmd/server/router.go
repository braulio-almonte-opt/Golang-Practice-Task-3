package main

import (
	"github.com/braulio-almonte-opt/task3-moviesAPI/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(actorController *controllers.Ac) *gin.Engine {

	r := gin.Default()

	api := r.Group("/api/v1")
	{

		actors := api.Group("/actors")
		{
			actors.POST("", actorController.CreateActor)
			actors.GET("", actorController.GetAllActors)
			actors.GET("/:id", actorController.GetActorByID)
			actors.PUT("/:id", actorController.UpdateActor)
			actors.DELETE("/:id", actorController.DeleteActor)
		}
	}

	return r
}
