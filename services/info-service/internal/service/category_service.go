package service

import (
	"Info_Service/internal/model"
	"Info_Service/internal/repository"
	"context"
)

type CategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.SpeciesRepository) *SpeciesService {
	return &SpeciesService{repo: repo}
}

func (s *CategoryService) GetAll(ctx context.Context) ([]model.Category, error) {
	return s.repo.GetAll(ctx)
}
