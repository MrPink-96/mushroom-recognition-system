package main

import (
	"api-gateway/internal/client"
	"api-gateway/internal/config"
	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"
	"api-gateway/internal/service"
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()

	httpClient := &http.Client{Timeout: 10 * time.Second}

	mlClient := client.NewMLClient(cfg.MLURL, httpClient)
	infoClient := client.NewInfoClient(cfg.InfoURL, httpClient)

	predictService := service.NewPredictService(mlClient, infoClient, cfg.ImageURL)

	healthHandler := handler.NewHealthHandler(mlClient, infoClient)
	predictHandler := handler.NewPredictHandler(predictService)
	speciesHandler := handler.NewSpeciesHandler(infoClient)
	r := gin.Default()

	r.Use(cors.Default())
	r.Use(middleware.TimeoutMiddleware(5 * time.Second))

	r.GET("/health", healthHandler.Check)
	r.POST("/predict", predictHandler.Predict)
	r.GET("/species", speciesHandler.GetAll)
	r.GET("/species/filter", speciesHandler.GetFiltered)

	srv := http.Server{
		Addr:         cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("API Gateway started on port %s", cfg.Port)
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
