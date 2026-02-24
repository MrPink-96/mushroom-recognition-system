package repository

import (
	"Info_Service/internal/model"
	"context"
	"github.com/jmoiron/sqlx"
)

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]model.Category, error)
}

type categoryRepo struct {
	db *sqlx.DB
}

const getAllQuery = `SELECT id, name, description FROM categories`

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) GetAll(ctx context.Context) ([]model.Category, error) {
	var result []model.Category
	err := r.db.SelectContext(ctx, &result, getAllQuery)
	return result, err
}
