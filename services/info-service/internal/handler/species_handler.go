package handler

import (
	"Info_Service/internal/dto"
	appErr "Info_Service/internal/errors"
	"Info_Service/internal/service"
	"errors"
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

func parseID(idStr string) (int64, error) {
	if idStr == "" {
		return 0, appErr.ErrInvalidID
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, appErr.ErrInvalidID
	}

	if id <= 0 {
		return 0, appErr.ErrInvalidID
	}

	return id, nil
}

func parsePageAndLimit(pageStr, limitStr string) (int, int, error) {
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, 0, appErr.ErrInvalidPage
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, appErr.ErrInvalidLimit
	}

	return page, limit, nil
}

func (h *SpeciesHandler) GetAll(c *gin.Context) {
	page, limit, err := parsePageAndLimit(c.DefaultQuery("page", "1"), c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	result, err := h.service.GetAll(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.ErrInternal})
		return
	}
	c.JSON(http.StatusOK, result)

}

func (h *SpeciesHandler) GetByID(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	result, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, appErr.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": appErr.ErrNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.ErrInternal})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SpeciesHandler) GetByCategory(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	page, limit, err := parsePageAndLimit(c.DefaultQuery("page", "1"), c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	result, err := h.service.GetByCategory(c.Request.Context(), id, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.ErrInternal})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *SpeciesHandler) SearchByName(c *gin.Context) {
	page, limit, err := parsePageAndLimit(c.DefaultQuery("page", "1"), c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	name := c.Query("name")
	result, err := h.service.SearchByName(c.Request.Context(), name, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.ErrInternal})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *SpeciesHandler) GetFiltered(c *gin.Context) {
	page, limit, err := parsePageAndLimit(c.DefaultQuery("page", "1"), c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	filter := dto.SpeciesFilter{Page: page, Limit: limit}
	if name := c.Query("name"); name != "" {
		filter.Name = &name
	}

	if category := c.Query("category"); category != "" {
		id, err := strconv.ParseInt(category, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": appErr.ErrInvalidCategoryID})
			return
		}
		filter.CategoryID = &id
	}

	if edibility := c.Query("edibility"); edibility != "" {
		val, err := strconv.Atoi(edibility)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": appErr.ErrInvalidEdibility})
			return
		}
		filter.Edibility = &val

	}

	if toxicity := c.Query("toxicity_max"); toxicity != "" {
		val, err := strconv.Atoi(toxicity)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": appErr.ErrInvalidToxicity})
			return
		}
		filter.ToxicityMax = &val

	}

}
