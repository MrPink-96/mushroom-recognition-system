package handler

import (
	"Info_Service/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type SpeciesHandler struct {
	service *service.SpeciesService
}

func NewSpeciesHandler(service *service.SpeciesService) *SpeciesHandler {
	return &SpeciesHandler{service: service}
}

func (h *SpeciesHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	species, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, species)
}
