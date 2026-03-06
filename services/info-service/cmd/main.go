package main

import (
	"Info_Service/internal/config"
	"Info_Service/internal/db"
	"Info_Service/internal/handler"
	"Info_Service/internal/repository"
	"Info_Service/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	cfg := config.Load()

	database, err := db.NewPostgres(cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}

	speciesRepo := repository.NewSpeciesRepository(database)
	speciesService := service.NewSpeciesService(speciesRepo)
	speciesHandler := handler.NewSpeciesHandler(speciesService)
	categoryRepo := repository.NewCategoryRepository(database)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	healthHandler := handler.NewHealthHandler(database)

	r := gin.Default()
	r.GET("/health", healthHandler.Check)

	r.GET("/species/", speciesHandler.GetAll)
	r.GET("/species/:id", speciesHandler.GetByID)
	r.GET("/species/search", speciesHandler.SearchByName)
	r.GET("/species/category/:id", speciesHandler.GetByCategory)
	r.GET("/categories", categoryHandler.GetAll)
	r.GET("/categories/:id", categoryHandler.GetByID)
	r.GET("/categories/search", categoryHandler.SearchByName)
	r.GET("/species/filter", speciesHandler.GetFiltered)

	log.Println("Info Servicetarted on: 8080")
	r.Run(":" + cfg.Port)
}
