package handler

import (
	"api-gateway/internal/client"
	appErr "api-gateway/internal/errors"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthHandler struct {
	mlClient   client.MLClient
	infoClient client.InfoClient
}

func NewHealthHandler(mlClient client.MLClient, infoClient client.InfoClient) *HealthHandler {
	return &HealthHandler{mlClient: mlClient, infoClient: infoClient}
}

func (h *HealthHandler) Check(c *gin.Context) {
	err := h.mlClient.Health(c.Request.Context())
	if err != nil {
		if errors.Is(err, appErr.ErrInfoUnavailable) {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "ml-service down",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": appErr.ErrInternal.Error(),
		})
		return
	}
	err = h.infoClient.Health(c.Request.Context())
	if err != nil {
		if errors.Is(err, appErr.ErrInfoUnavailable) {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "info-service down",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": appErr.ErrInternal.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"service":      "api-gateway",
		"status":       "ok",
		"info-service": "up",
		"ml-service":   "up",
	})
}
