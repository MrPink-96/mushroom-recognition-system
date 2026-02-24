package service

import (
	"Info_Service/internal/dto"
	"Info_Service/internal/repository"
	"context"
)

type SpeciesService struct {
	repo repository.SpeciesRepository
}

func NewSpeciesService(repo repository.SpeciesRepository) *SpeciesService {
	return &SpeciesService{repo: repo}
}

func normalize(page, limit int) (int, int) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	return page, limit
}

func (s *SpeciesService) GetByID(ctx context.Context, id int64) (*dto.SpeciesResponse, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SpeciesService) GetAll(ctx context.Context, page, limit int) ([]dto.SpeciesResponse, error) {
	page, limit = normalize(page, limit)
	return s.repo.GetAll(ctx, page, limit)
}

func (s *SpeciesService) GetByCategory(ctx context.Context, categoryID int64, page, limit int) ([]dto.SpeciesResponse, error) {
	page, limit = normalize(page, limit)
	return s.repo.GetByCategory(ctx, categoryID, page, limit)
}

func (s *SpeciesService) SearchByName(ctx context.Context, name string, page, limit int) ([]dto.SpeciesResponse, error) {
	page, limit = normalize(page, limit)
	return s.repo.SearchByName(ctx, name, page, limit)
}
