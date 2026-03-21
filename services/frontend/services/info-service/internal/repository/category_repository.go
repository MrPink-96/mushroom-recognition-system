package repository

import (
	"Info_Service/internal/dto"
	appErr "Info_Service/internal/errors"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type CategoryRepository interface {
	GetByID(ctx context.Context, id int64) (*dto.CategoryResponse, error)
	GetAll(ctx context.Context, page, limit int, sortBy, order string) ([]dto.CategoryResponse, int, error)
	SearchByName(ctx context.Context, name string, page, limit int) ([]dto.CategoryResponse, int, error)
}

type categoryRepo struct {
	db *sqlx.DB
}

const categoryBaseQuery = `SELECT id, name, description FROM categories c`

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) GetByID(ctx context.Context, id int64) (*dto.CategoryResponse, error) {
	query := categoryBaseQuery + ` WHERE c.id=$1`
	var result dto.CategoryResponse
	err := r.db.GetContext(ctx, &result, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErr.ErrNotFound
		}
		return nil, err
	}
	return &result, nil

}

func (r *categoryRepo) GetAll(ctx context.Context, page, limit int, sortField, sortOrder string) ([]dto.CategoryResponse, int, error) {
	offset := (page - 1) * limit

	countQuery := `SELECT COUNT(*) FROM categories`
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`%s ORDER BY c.%s %s LIMIT $1 OFFSET $2`,
		categoryBaseQuery, sortField, sortOrder)
	var result []dto.CategoryResponse
	err := r.db.SelectContext(ctx, &result, query, limit, offset)
	return result, total, err
}

func (r *categoryRepo) SearchByName(ctx context.Context, name string, page, limit int) ([]dto.CategoryResponse, int, error) {
	offset := (page - 1) * limit

	patern := "%" + name + "%"
	countQuery := `SELECT COUNT(*) FROM categories WHERE name ILIKE $1`
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, patern); err != nil {
		return nil, 0, err
	}

	query := categoryBaseQuery + ` WHERE c.name ILIKE $1 ORDER BY c.id ASC LIMIT $2 OFFSET $3`
	var result []dto.CategoryResponse
	err := r.db.SelectContext(ctx, &result, query, patern, limit, offset)
	return result, total, err
}
