package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/braulio-almonte-opt/task3-moviesAPI/controllers"
	"github.com/braulio-almonte-opt/task3-moviesAPI/dao"
	"github.com/braulio-almonte-opt/task3-moviesAPI/models"
	"github.com/braulio-almonte-opt/task3-moviesAPI/service"
)

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {
	client, err := models.ConnectDatabase()
	if err != nil {
		return err
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("disconnect MongoDB: %v", err)
		}
	}()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	dbName := "moviesAPI"
	actorDAO := dao.NewActorUserDAO(client, dbName)
	actorService := service.NewActorService(actorDAO)
	actorController := controllers.NewActorController(actorService)

	r, err := SetupRouter(actorController)
	if err != nil {
		return fmt.Errorf("initialize application dependencies: %w", err)
	}

	log.Println("Server started on :8080")
	if err := r.Run(":8080"); err != nil {
		return err
	}
	return nil
}
