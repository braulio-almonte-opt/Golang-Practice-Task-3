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

		films := api.Group("/films")
		{
			films.POST("", filmController.CreateFilm)
			films.GET("", filmController.GetAllFilms)
			films.GET("/:id", filmController.GetFilmByID)
			films.PUT("/:id", filmController.UpdateFilm)
			films.DELETE("/:id", filmController.DeleteFilm)

			films.GET("/director/:directorId", filmController.GetFilmsByDirector)
			films.GET("/actor/:actorId", filmController.GetFilmsByActor)
		}

	}

	return r, nil
}
