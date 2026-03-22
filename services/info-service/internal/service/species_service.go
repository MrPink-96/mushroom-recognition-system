package service

import (
	"Info_Service/internal/config"
	"Info_Service/internal/dto"
	"Info_Service/internal/repository"
	"context"
	"fmt"
	"math"
	"strings"
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

func (s *SpeciesService) attachImages(ctx context.Context, species []dto.SpeciesResponse) error {

	if len(species) == 0 {
		return nil
	}

	ids := make([]int64, 0, len(species))

	for _, sp := range species {
		ids = append(ids, sp.ID)
	}

	imagesMap, err := s.repo.GetImagesBySpeciesIDs(ctx, ids)
	if err != nil {
		return err
	}
	//
	cfg := config.Load()
	imageBaseURL := fmt.Sprintf("http://localhost:%s", cfg.Port)
	for i := range species {
		if imgs, ok := imagesMap[species[i].ID]; ok {

			urls := make([]string, len(imgs))
			for j, img := range imgs {
				if strings.HasPrefix(img, "http") {
					urls[j] = img
				} else {
					// Иначе добавляем базовый URL
					urls[j] = fmt.Sprintf("%s%s", imageBaseURL, img)
				}
			}
			species[i].Images = urls
		} else {
			species[i].Images = []string{}
		}
	}
	//
	/*
		for i := range species {
			if imgs, ok := imagesMap[species[i].ID]; ok {
				species[i].Images = imgs
			} else {
				species[i].Images = []string{}
			}
		}
	*/
	return nil
}

func (s *SpeciesService) GetByID(ctx context.Context, id int64) (*dto.SpeciesResponse, error) {

	species, err := s.repo.GetSpeciesByID(ctx, id)
	if err != nil {
		return nil, err
	}

	list := []dto.SpeciesResponse{*species}

	if err := s.attachImages(ctx, list); err != nil {
		return nil, err
	}

	*species = list[0]

	return species, nil
}

func (s *SpeciesService) GetByIDs(ctx context.Context, ids []int64) ([]dto.SpeciesResponse, error) {

	if len(ids) == 0 {
		return []dto.SpeciesResponse{}, nil
	}

	data, err := s.repo.GetSpeciesByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	if err := s.attachImages(ctx, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *SpeciesService) GetAll(ctx context.Context, page, limit int, sortBy, order string) (dto.PaginatedSpeciesResponse, error) {
	page, limit = normalizePagination(page, limit)
	order = normalizeSortOrder(order)

	data, total, err := s.repo.GetAll(ctx, page, limit, sortBy, order)
	if err != nil {
		return dto.PaginatedSpeciesResponse{}, err
	}

	if err := s.attachImages(ctx, data); err != nil {
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

	if err := s.attachImages(ctx, data); err != nil {
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
	if err := s.attachImages(ctx, data); err != nil {
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
	filter.Order = normalizeSortOrder(filter.Order)

	data, total, err := s.repo.GetFiltered(ctx, filter)
	if err != nil {
		return dto.PaginatedSpeciesResponse{}, err
	}
	if err := s.attachImages(ctx, data); err != nil {
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
