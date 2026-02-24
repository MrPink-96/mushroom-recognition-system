package main

import (
	"Info_Service/internal/db"
	"Info_Service/internal/handler"
	"Info_Service/internal/repository"
	"Info_Service/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	dsn := "postgres://user:pass@localhost:5432/mushrooms?sslmode=disable"

	database, err := db.NewPostgres(dsn)
	if err != nil {
		panic(err)
	}

	speciesRepo := repository.NewSpeciesRepository(database)
	speciesService := service.NewSpeciesService(speciesRepo)
	speciesHandler := handler.NewSpeciesHandler(speciesService)

	r := gin.Default()

	r.GET("/species/:id", speciesHandler.GetByID)

	r.Run(":8080")
}
