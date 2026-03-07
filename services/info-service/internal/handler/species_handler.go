package handler

import (
	"Info_Service/internal/dto"
	appErr "Info_Service/internal/errors"
	"Info_Service/internal/service"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SpeciesHandler struct {
	service *service.SpeciesService
}

func NewSpeciesHandler(service *service.SpeciesService) *SpeciesHandler {
	return &SpeciesHandler{service: service}
}

func (h *SpeciesHandler) GetByID(c *gin.Context) {
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

func (h *SpeciesHandler) GetAll(c *gin.Context) {
	page, limit, err := parsePageAndLimit(c.DefaultQuery("page", "1"), c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sortBy := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "asc")

	result, err := h.service.GetAll(c.Request.Context(), page, limit, sortBy, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)

}

func (h *SpeciesHandler) GetByCategory(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page, limit, err := parsePageAndLimit(c.DefaultQuery("page", "1"), c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.GetByCategory(c.Request.Context(), id, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.ErrInternal.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *SpeciesHandler) SearchByName(c *gin.Context) {
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

func (h *SpeciesHandler) GetFiltered(c *gin.Context) {
	page, limit, err := parsePageAndLimit(c.DefaultQuery("page", "1"), c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := dto.SpeciesFilter{
		Page:  page,
		Limit: limit,
		Sort:  c.DefaultQuery("sort", "id"),
		Order: c.DefaultQuery("order", "asc"),
	}
	if name := c.Query("name"); name != "" {
		filter.Name = &name
	}

	if category := c.Query("category"); category != "" {
		id, err := strconv.ParseInt(category, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": appErr.ErrInvalidCategoryID.Error()})
			return
		}
		filter.CategoryID = &id
	}

	if edibility := c.Query("edibility"); edibility != "" {
		val, err := strconv.Atoi(edibility)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": appErr.ErrInvalidEdibility.Error()})
			return
		}
		filter.Edibility = &val

	}

	if toxicity := c.Query("toxicity_max"); toxicity != "" {
		val, err := strconv.Atoi(toxicity)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": appErr.ErrInvalidToxicity.Error()})
			return
		}
		filter.ToxicityMax = &val

	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	result, err := h.service.GetFiltered(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.ErrInternal.Error()})
		return
	}

	c.JSON(http.StatusOK, result)

}
