package service

import (
	"Info_Service/internal/model"
	"Info_Service/internal/repository"
	"context"
)

type SpeciesService struct {
	repo repository.SpeciesRepository
}

func NewSpeciesService(repo repository.SpeciesRepository) *SpeciesService {
	return &SpeciesService{repo: repo}
}

func (s *SpeciesService) GetByID(ctx context.Context, id int64) (*model.Species, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SpeciesService) GetAll(ctx context.Context) ([]model.Species, error) {
	return s.repo.GetAll(ctx)
}
