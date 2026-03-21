package repository

import (
	"Info_Service/internal/dto"
	appErr "Info_Service/internal/errors"
	"Info_Service/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strings"
)

type SpeciesRepository interface {
	GetAll(ctx context.Context, page, limit int, sortBy, order string) ([]dto.SpeciesResponse, int, error)
	GetSpeciesByID(ctx context.Context, id int64) (*dto.SpeciesResponse, error)
	GetSpeciesByIDs(ctx context.Context, ids []int64) ([]dto.SpeciesResponse, error)
	GetByCategory(ctx context.Context, categoryID int64, page, limit int) ([]dto.SpeciesResponse, int, error)
	SearchByName(ctx context.Context, name string, page, limit int) ([]dto.SpeciesResponse, int, error)
	GetFiltered(ctx context.Context, filter dto.SpeciesFilter) ([]dto.SpeciesResponse, int, error)
	GetImagesBySpeciesIDs(ctx context.Context, ids []int64) (map[int64][]string, error)
}

type speciesRepo struct {
	db *sqlx.DB
}

func NewSpeciesRepository(db *sqlx.DB) SpeciesRepository {
	return &speciesRepo{db: db}
}

const speciesBaseQuery = `
SELECT 
 	s.id AS "id",
    s.scientific_name AS "scientific_name",
    s.common_name AS "common_name",
    s.description AS "description",
    s.edibility AS "edibility",
    s.toxicity_level AS "toxicity_level",
    
	c.id   AS "category.id",
	c.name AS "category.name"
FROM species s
JOIN categories c ON s.category_id = c.id
`

var allowedSortFields = map[string]string{
	"id":              "s.id",
	"scientific_name": "s.scientific_name",
	"common_name":     "s.common_name",
	"toxicity":        "s.toxicity_level",
	"edibility":       "s.edibility",
}

func (r *speciesRepo) GetSpeciesByID(ctx context.Context, id int64) (*dto.SpeciesResponse, error) {
	query := speciesBaseQuery + ` WHERE s.id = $1`
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

func (r *speciesRepo) GetSpeciesByIDs(ctx context.Context, ids []int64) ([]dto.SpeciesResponse, error) {

	if len(ids) == 0 {
		return []dto.SpeciesResponse{}, nil
	}

	query := speciesBaseQuery + `
	WHERE s.id = ANY($1)
	ORDER BY array_position($1, s.id)
	`

	var result []dto.SpeciesResponse

	err := r.db.SelectContext(ctx, &result, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *speciesRepo) GetAll(ctx context.Context, page, limit int, sortBy, order string) ([]dto.SpeciesResponse, int, error) {
	offset := (page - 1) * limit

	var total int
	countQuery := `SELECT COUNT(*) FROM species`
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, err
	}

	sortField, ok := allowedSortFields[sortBy]
	if !ok {
		sortField = "s.id"
	}

	query := fmt.Sprintf(`%s ORDER BY %s %s LIMIT $1 OFFSET $2`,
		speciesBaseQuery, sortField, order)

	var result []dto.SpeciesResponse
	err := r.db.SelectContext(ctx, &result, query, limit, offset)
	return result, total, err
}

func (r *speciesRepo) GetByCategory(ctx context.Context, categoryID int64, page, limit int) ([]dto.SpeciesResponse, int, error) {
	offset := (page - 1) * limit

	var total int
	countQuery := `SELECT COUNT(*) FROM species s WHERE s.category_id = $1`
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, err
	}

	query := speciesBaseQuery + ` WHERE s.category_id = $1 ORDER BY s.id ASC LIMIT $2 OFFSET $3`
	var result []dto.SpeciesResponse
	err := r.db.SelectContext(ctx, &result, query, categoryID, limit, offset)
	return result, total, err
}

func (r *speciesRepo) SearchByName(ctx context.Context, name string, page, limit int) ([]dto.SpeciesResponse, int, error) {
	offset := (page - 1) * limit

	var total int
	pattern := "%" + name + "%"
	countQuery := `SELECT COUNT(*)	FROM species s
				WHERE s.common_name ILIKE $1 OR s.scientific_name ILIKE $1`
	if err := r.db.GetContext(ctx, &total, countQuery, pattern); err != nil {
		return nil, 0, err
	}

	query := speciesBaseQuery + ` WHERE s.common_name ILIKE $1 OR s.scientific_name ILIKE $1
			ORDER BY s.id ASC LIMIT $2 OFFSET $3`

	var result []dto.SpeciesResponse
	err := r.db.SelectContext(ctx, &result, query, pattern, limit, offset)
	return result, total, err
}

func (r *speciesRepo) GetFiltered(
	ctx context.Context,
	filter dto.SpeciesFilter,
) ([]dto.SpeciesResponse, int, error) {
	if filter.Cursor != nil {
		filter.Sort = "id"
		filter.Order = "asc"
	}

	where := []string{}
	args := []interface{}{}
	argPos := 1

	if filter.Name != nil {
		where = append(where,
			fmt.Sprintf("(s.common_name ILIKE $%d OR s.scientific_name ILIKE $%d)", argPos, argPos))
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
		where = append(where, fmt.Sprintf("s.toxicity_level <= $%d", argPos))
		args = append(args, *filter.ToxicityMax)
		argPos++
	}

	if filter.Cursor != nil {
		where = append(where, fmt.Sprintf("s.id > $%d", argPos))
		args = append(args, *filter.Cursor)
		argPos++
	}

	query := speciesBaseQuery
	countQuery := `SELECT COUNT(*) FROM species s`

	if len(where) > 0 {
		condition := " WHERE " + strings.Join(where, " AND ")
		query += condition
		countQuery += condition
	}

	var total int
	if filter.Cursor == nil {
		if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
			return nil, 0, err
		}
	}

	sortField, ok := allowedSortFields[filter.Sort]
	if !ok {
		sortField = "s.id"
	}

	query += fmt.Sprintf(" ORDER BY %s %s, s.id ASC LIMIT $%d", sortField, filter.Order, argPos)
	args = append(args, filter.Limit)

	if filter.Cursor == nil {
		offset := (filter.Page - 1) * filter.Limit
		query += fmt.Sprintf(" OFFSET $%d", argPos+1)
		args = append(args, offset)
	}

	var result []dto.SpeciesResponse
	err := r.db.SelectContext(ctx, &result, query, args...)

	return result, total, err
}
func (r *speciesRepo) GetImagesBySpeciesIDs(ctx context.Context, ids []int64) (map[int64][]string, error) {

	if len(ids) == 0 {
		return map[int64][]string{}, nil
	}

	query := `
		SELECT species_id, image_path
		FROM species_images
		WHERE species_id = ANY($1)
		ORDER BY species_id, is_primary DESC, id ASC
	`

	var rows []model.SpeciesImage

	err := r.db.SelectContext(ctx, &rows, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	result := make(map[int64][]string, len(ids))

	for _, row := range rows {
		result[row.SpeciesID] = append(result[row.SpeciesID], row.ImagePath)
	}

	return result, nil
}
