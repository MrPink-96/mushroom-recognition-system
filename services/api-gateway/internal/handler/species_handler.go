package handler

import (
	"api-gateway/internal/client"
	"github.com/gin-gonic/gin"
)

type SpeciesHandler struct {
	info client.InfoClient
}

func NewSpeciesHandler(info client.InfoClient) *SpeciesHandler {
	return &SpeciesHandler{info: info}
}

func (h *SpeciesHandler) GetAll(c *gin.Context) {

}

func (h *SpeciesHandler) GetFiltered(c *gin.Context) {

}
