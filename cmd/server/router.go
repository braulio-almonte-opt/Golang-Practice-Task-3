package main

import (
	"github.com/braulio-almonte-opt/task3-moviesAPI/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(actorController *controllers.ActorController, directorController *controllers.DirectorController) (*gin.Engine, error) {

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

		directors := api.Group("/directors")
		{
			directors.POST("", directorController.CreateDirector)
			directors.GET("", directorController.GetAllDirectors)
			directors.GET("/:id", directorController.GetDirectorByID)
			directors.PUT("/:id", directorController.UpdateDirector)
			directors.DELETE("/:id", directorController.DeleteDirector)

		}
	}

	return r, nil
}
