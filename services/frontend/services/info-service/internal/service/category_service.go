package service

import (
	"Info_Service/internal/dto"
	"Info_Service/internal/repository"
	"context"
	"math"
	"strings"
)

type CategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) buildPaginated(data []dto.CategoryResponse, total, page, limit int) dto.PaginatedCategoryResponse {
	pages := int(math.Ceil(float64(total) / float64(limit)))

	return dto.PaginatedCategoryResponse{
		Data: data,
		Meta: dto.Meta{
			Page:  page,
			Limit: limit,
			Total: total,
			Pages: pages,
		},
	}
}

func (s *CategoryService) normalizeSortField(sortBy string) string {
	if strings.ToLower(sortBy) == "name" {
		return "name"
	}
	return "id"
}

func (s *CategoryService) GetByID(ctx context.Context, id int64) (*dto.CategoryResponse, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) GetAll(ctx context.Context, page, limit int, sortField, sortOrder string) (dto.PaginatedCategoryResponse, error) {
	page, limit = normalizePagination(page, limit)
	sortField = s.normalizeSortField(sortField)
	sortOrder = normalizeSortOrder(sortOrder)

	data, total, err := s.repo.GetAll(ctx, page, limit, sortField, sortOrder)
	if err != nil {
		return dto.PaginatedCategoryResponse{}, err
	}

	return s.buildPaginated(data, total, page, limit), nil
}

func (s *CategoryService) SearchByName(ctx context.Context, name string, page, limit int) (dto.PaginatedCategoryResponse, error) {
	page, limit = normalizePagination(page, limit)

	data, total, err := s.repo.SearchByName(ctx, name, page, limit)
	if err != nil {
		return dto.PaginatedCategoryResponse{}, err
	}
	return s.buildPaginated(data, total, page, limit), nil
}
