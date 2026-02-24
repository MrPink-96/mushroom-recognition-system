package repository

import (
	"Info_Service/internal/dto"
	appErr "Info_Service/internal/errors"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type SpeciesRepository interface {
	GetByID(ctx context.Context, id int64) (*dto.SpeciesResponse, error)
	GetAll(ctx context.Context, page, limit int) ([]dto.SpeciesResponse, error)
	GetByCategory(ctx context.Context, categoryID int64, page, limit int) ([]dto.SpeciesResponse, error)
	SearchByName(ctx context.Context, name string, page, limit int) ([]dto.SpeciesResponse, error)
	GetFiltered(ctx context.Context, edibility *int, toxicityMax *int, limit, offset int) ([]dto.SpeciesResponse, error)
}

type speciesRepo struct {
	db *sqlx.DB
}

func NewSpeciesRepository(db *sqlx.DB) SpeciesRepository {
	return &speciesRepo{db: db}
}

const baseQuery = `
SELECT 
	s.id,
	s.scientific_name,
	s.common_name,
	s.description,
	s.edibility,
	s.toxicity_level,
	s.reference_image_url,
	c.name AS category_name
FROM species s
JOIN categories c ON s.category_id = c.id
`

func (r *speciesRepo) GetByID(ctx context.Context, id int64) (*dto.SpeciesResponse, error) {
	query := baseQuery + " WHERE s.id = $1"
	var result dto.SpeciesResponse
	err := r.db.GetContext(ctx, &result, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErr.ErrNotFound
		}
		return nil, err
	}
	return &result, nil
}

func (r *speciesRepo) GetAll(ctx context.Context, page, limit int) ([]dto.SpeciesResponse, error) {
	offset := (page - 1) * limit
	query := baseQuery + "LIMIT $1 OFFSET $2"
	var result []dto.SpeciesResponse

	err := r.db.SelectContext(ctx, &result, query, limit, offset)
	return result, err
}

func (r *speciesRepo) GetByCategory(ctx context.Context, categoryID int64, page, limit int) ([]dto.SpeciesResponse, error) {
	offset := (page - 1) * limit
	query := baseQuery + "WHERE s.category_id = $1 LIMIT $2 OFFSET $3"

	var result []dto.SpeciesResponse
	err := r.db.SelectContext(ctx, &result, query, categoryID, limit, offset)
	return result, err
}

func (r *speciesRepo) SearchByName(ctx context.Context, name string, page, limit int) ([]dto.SpeciesResponse, error) {
	offset := (page - 1) * limit
	query := baseQuery + `WHERE s.common_name ILIKE '%' || $1 || '%'  OR s.scientific_name ILIKE '%' || $1 || '%'
							LIMIT $2 OFFSET $3`

	var result []dto.SpeciesResponse
	err := r.db.SelectContext(ctx, &result, query, name, limit, offset)
	return result, err
}

func (r *speciesRepo) GetFiltered(ctx context.Context, edibility *int, toxicityMax *int, limit, offset int) ([]dto.SpeciesResponse, error) {

}
