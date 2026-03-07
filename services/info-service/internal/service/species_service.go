package service

import (
	"Info_Service/internal/dto"
	"Info_Service/internal/repository"
	"context"
	"math"
)

type SpeciesService struct {
	repo repository.SpeciesRepository
}

func NewSpeciesService(repo repository.SpeciesRepository) *SpeciesService {
	return &SpeciesService{repo: repo}
}

func (s *SpeciesService) buildPaginated(data []dto.SpeciesResponse, total, page, limit int) dto.PaginatedSpeciesResponse {
	pages := int(math.Ceil(float64(total) / float64(limit)))

	return dto.PaginatedSpeciesResponse{
		Data: data,
		Meta: dto.Meta{
			Page:  page,
			Limit: limit,
			Total: total,
			Pages: pages,
		},
	}
}

func (s *SpeciesService) GetByID(ctx context.Context, id int64) (*dto.SpeciesResponse, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SpeciesService) GetAll(ctx context.Context, page, limit int, sortBy, order string) (dto.PaginatedSpeciesResponse, error) {
	page, limit = normalizePagination(page, limit)

	data, total, err := s.repo.GetAll(ctx, page, limit, sortBy, order)
	if err != nil {
		return dto.PaginatedSpeciesResponse{}, err
	}

	return s.buildPaginated(data, total, page, limit), nil
}

func (s *SpeciesService) GetByCategory(ctx context.Context, categoryID int64, page, limit int) (dto.PaginatedSpeciesResponse, error) {
	page, limit = normalizePagination(page, limit)

	data, total, err := s.repo.GetByCategory(ctx, categoryID, page, limit)
	if err != nil {
		return dto.PaginatedSpeciesResponse{}, err
	}

	return s.buildPaginated(data, total, page, limit), nil

}

func (s *SpeciesService) SearchByName(ctx context.Context, name string, page, limit int) (dto.PaginatedSpeciesResponse, error) {
	page, limit = normalizePagination(page, limit)

	data, total, err := s.repo.SearchByName(ctx, name, page, limit)
	if err != nil {
		return dto.PaginatedSpeciesResponse{}, err
	}
	return s.buildPaginated(data, total, page, limit), nil
}

func (s *SpeciesService) GetFiltered(ctx context.Context, filter dto.SpeciesFilter) (dto.PaginatedSpeciesResponse, error) {
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 10
	}

	if filter.Cursor == nil {
		if filter.Page <= 0 {
			filter.Page = 1
		}
	}

	data, total, err := s.repo.GetFiltered(ctx, filter)
	if err != nil {
		return dto.PaginatedSpeciesResponse{}, err
	}

	meta := dto.Meta{
		Limit: filter.Limit,
	}

	if filter.Cursor == nil {
		pages := (total + filter.Limit - 1) / filter.Limit
		meta.Page = filter.Page
		meta.Total = total
		meta.Pages = pages
	}

	if filter.Cursor != nil && len(data) > 0 {
		lastID := data[len(data)-1].ID
		meta.NextCursor = &lastID
	}

	return dto.PaginatedSpeciesResponse{
		Data: data,
		Meta: meta,
	}, nil
}
