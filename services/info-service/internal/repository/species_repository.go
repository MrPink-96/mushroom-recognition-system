package repository

import (
	"Info_Service/internal/dto"
	appErr "Info_Service/internal/errors"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type SpeciesRepository interface {
	GetByID(ctx context.Context, id int64) (*dto.SpeciesResponse, error)
	GetAll(ctx context.Context, page, limit int) ([]dto.SpeciesResponse, int, error)
	GetByCategory(ctx context.Context, categoryID int64, page, limit int) ([]dto.SpeciesResponse, int, error)
	SearchByName(ctx context.Context, name string, page, limit int) ([]dto.SpeciesResponse, int, error)
	GetFiltered(ctx context.Context, filter dto.SpeciesFilter) ([]dto.SpeciesResponse, int, error)
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
    c.id   AS category_id,
    c.name AS category_name
FROM species s
JOIN categories c ON s.category_id = c.id
`

func (r *speciesRepo) GetByID(ctx context.Context, id int64) (*dto.SpeciesResponse, error) {
	query := baseQuery + ` WHERE s.id = $1`
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

func (r *speciesRepo) GetAll(ctx context.Context, page, limit int) ([]dto.SpeciesResponse, int, error) {
	offset := (page - 1) * limit

	var total int
	countQuery := `SELECT COUNT(*) FROM species`
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, err
	}

	query := baseQuery + ` ORDER BY s.id ASC LIMIT $1 OFFSET $2`
	var result []dto.SpeciesResponse
	err := r.db.SelectContext(ctx, &result, query, limit, offset)
	return result, total, err
}

func (r *speciesRepo) GetByCategory(ctx context.Context, categoryID int64, page, limit int) ([]dto.SpeciesResponse, int, error) {
	offset := (page - 1) * limit

	var total int
	countQuery := `SELECT Count(*) FROM species WHERE s.category_id = $1`
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, err
	}

	query := baseQuery + ` WHERE s.category_id = $1 ORDER BY s.id ASC LIMIT $2 OFFSET $3`
	var result []dto.SpeciesResponse
	err := r.db.SelectContext(ctx, &result, query, categoryID, limit, offset)
	return result, total, err
}

func (r *speciesRepo) SearchByName(ctx context.Context, name string, page, limit int) ([]dto.SpeciesResponse, int, error) {
	offset := (page - 1) * limit

	var total int
	pattern := "%" + name + "%"
	countQuery := `SELECT COUNT(*)	FROM species
				WHERE common_name ILIKE $1 OR scientific_name ILIKE $1`
	if err := r.db.GetContext(ctx, &total, countQuery, pattern); err != nil {
		return nil, 0, err
	}

	query := baseQuery + ` WHERE s.common_name ILIKE $1 OR s.scientific_name ILIKE $1
			ORDER BY s.id ASC LIMIT $2 OFFSET $3`

	var result []dto.SpeciesResponse
	err := r.db.SelectContext(ctx, &result, query, pattern, limit, offset)
	return result, total, err
}

func (r *speciesRepo) GetFiltered(ctx context.Context, filter dto.SpeciesFilter) ([]dto.SpeciesResponse, int, error) {
	where := []string{}
	args := []interface{}{}
	argPos := 1

	if filter.Name != nil {
		where = append(where, fmt.Sprintf("(s.common_name ILIKE $%d OR s.scientific_name ILIKE $%d)", argPos, argPos))
		args = append(args, "%"+*filter.Name+"%")
		argPos++
	}

	if filter.CategoryID != nil {
		where = append(where, fmt.Sprintf("s.category_id = $%d", argPos))
		args = append(args, *filter.CategoryID)
		argPos++
	}

	if filter.Edibility != nil {
		where = append(where, fmt.Sprintf("s.edibility = $%d", argPos))
		args = append(args, *filter.Edibility)
		argPos++
	}

	if filter.ToxicityMax != nil {
		where = append(where, fmt.Sprintf("s.edibility = $%d", argPos))
		args = append(args, *filter.ToxicityMax)
		argPos++
	}

	query := baseQuery
	countQuery := `SELECT COUNT(*) FROM species s`
	if len(where) > 0 {
		condition := " WHERE " + strings.Join(where, " AND ")
		query += condition
		countQuery += condition
	}

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	query += fmt.Sprintf(" ORDER BY s.id ASC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, filter.Limit, offset)

	var result []dto.SpeciesResponse
	err := r.db.SelectContext(ctx, &result, query, args...)

	return result, total, err

}
