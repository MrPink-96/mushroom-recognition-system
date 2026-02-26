package main

import (
	"Info_Service/internal/db"
	"Info_Service/internal/handler"
	"Info_Service/internal/repository"
	"Info_Service/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN not set")
	}

	database, err := db.NewPostgres(dsn)
	if err != nil {
		log.Fatal(err)
	}

	speciesRepo := repository.NewSpeciesRepository(database)
	speciesService := service.NewSpeciesService(speciesRepo)
	speciesHandler := handler.NewSpeciesHandler(speciesService)
	categoryRepo := repository.NewCategoryRepository(database)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	r := gin.Default()

	r.GET("/species/", speciesHandler.GetAll)
	r.GET("/species/:id", speciesHandler.GetByID)
	r.GET("/species/search", speciesHandler.SearchByName)
	r.GET("/species/category/:id", speciesHandler.GetByCategory)
	r.GET("/categories", categoryHandler.GetAll)
	r.GET("/species/filter", speciesHandler.GetFiltered)

	log.Println("Info Servicetarted on: 8080")
	r.Run(":8080")
}
