package handler

import (
	appErr "Info_Service/internal/errors"
	"Info_Service/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, appErr.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": appErr.ErrNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.ErrInternal.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	page, limit, err := parsePageAndLimit(c.DefaultQuery("page", "1"), c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sortBy := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "asc")

	result, err := h.service.GetAll(c.Request.Context(), page, limit, sortBy, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *CategoryHandler) SearchByName(c *gin.Context) {
	page, limit, err := parsePageAndLimit(c.DefaultQuery("page", "1"), c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := c.Query("name")
	if strings.TrimSpace(name) == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.ErrEmptySearchQuery.Error()})
		return
	}

	result, err := h.service.SearchByName(c.Request.Context(), name, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.ErrInternal.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
