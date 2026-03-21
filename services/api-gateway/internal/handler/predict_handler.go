package handler

import (
	appErr "api-gateway/internal/errors"
	"api-gateway/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

const maxFileSize = 10 << 20

type PredictHandler struct {
	service *service.PredictService
}

func NewPredictHandler(service *service.PredictService) *PredictHandler {
	return &PredictHandler{service: service}
}

func (h *PredictHandler) Predict(c *gin.Context) {
	if c.Request.ContentLength > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": appErr.ErrFileLarge.Error()})
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxFileSize)
	resp, err := h.service.Predict(c.Request.Context(), c.Request)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)

}
