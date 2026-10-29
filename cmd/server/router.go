package main

import (
	"github.com/braulio-almonte-opt/task3-moviesAPI/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	actorController *controllers.ActorController,
	directorController *controllers.DirectorController,
	filmController *controllers.FilmController) (*gin.Engine, error) {

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		actorController.RegisterRoutes(api)
		directorController.RegisterRoutes(api)
		filmController.RegisterRoutes(api)
	}

	return r, nil
}
