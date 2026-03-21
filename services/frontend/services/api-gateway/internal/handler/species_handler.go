package handler

import (
	"api-gateway/internal/client"
	appErr "api-gateway/internal/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

type SpeciesHandler struct {
	info client.InfoClient
}

func NewSpeciesHandler(info client.InfoClient) *SpeciesHandler {
	return &SpeciesHandler{info: info}
}

func (h *SpeciesHandler) proxy(c *gin.Context, path string) {
	query := url.Values{}

	for k, v := range c.Request.URL.Query() {
		for _, val := range v {
			query.Add(k, val)
		}
	}

	body, status, err := h.info.ProxyGet(c.Request.Context(), path, query)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": appErr.ErrInfoUnavailable.Error()})
		return
	}

	c.Data(status, "application/json", body)
}

func (h *SpeciesHandler) GetAll(c *gin.Context) {
	h.proxy(c, "/species")
}

func (h *SpeciesHandler) GetFiltered(c *gin.Context) {
	h.proxy(c, "/species/filter")
}
