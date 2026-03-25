package main

import (
	"Info_Service/internal/config"
	"Info_Service/internal/db"
	"Info_Service/internal/handler"
	"Info_Service/internal/middleware"
	"Info_Service/internal/repository"
	"Info_Service/internal/service"
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()
	if cfg.DSN == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	database, err := db.NewPostgres(cfg.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	speciesRepo := repository.NewSpeciesRepository(database)
	speciesService := service.NewSpeciesService(speciesRepo)
	speciesHandler := handler.NewSpeciesHandler(speciesService)

	categoryRepo := repository.NewCategoryRepository(database)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	healthHandler := handler.NewHealthHandler(database)

	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(cors.Default())
	r.Use(middleware.TimeoutMiddleware(5 * time.Second))

	//r.Static("/images", "./images")

	r.GET("/health", healthHandler.Check)

	r.GET("/species", speciesHandler.GetAll)
	r.GET("/species/:id", speciesHandler.GetByID)
	r.GET("/species/batch", speciesHandler.GetByIDs)
	r.GET("/species/search", speciesHandler.SearchByName)
	r.GET("/species/category/:id", speciesHandler.GetByCategory)
	r.GET("/species/filter", speciesHandler.GetFiltered)

	r.GET("/categories", categoryHandler.GetAll)
	r.GET("/categories/:id", categoryHandler.GetByID)
	r.GET("/categories/search", categoryHandler.SearchByName)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Info Service started on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
