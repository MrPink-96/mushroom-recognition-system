package handler

import (
	appErr "api-gateway/internal/errors"
	"api-gateway/internal/service"
	"github.com/gin-gonic/gin"
	"io"
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
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": appErr.ErrFileRequired.Error()})
		return
	}
	defer file.Close()

	if header.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": appErr.ErrFileLarge.Error()})
		return
	}
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": appErr.ErrInvalidFile.Error()})
		return
	}

	contentType := http.DetectContentType(buffer)
	if contentType != "image/jpeg" && contentType != "image/png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": appErr.ErrInvalidFileType.Error()})
		return
	}

	file.Seek(0, 0)

	image, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.ErrReadFile.Error()})
		return
	}

	resp, err := h.service.Predict(c.Request.Context(), image)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)

}
