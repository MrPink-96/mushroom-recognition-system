package repository

import (
	"Info_Service/internal/model"
	"context"
	"database/sql"
)

type SpeciesRepository interface {
	GetByID(ctx context.Context, id int64) (*model.Species, error)
	GetAll(ctx context.Context) ([]model.Species, error)
	GetByCategory(ctx context.Context, categoryID int64) ([]model.Species, error)
	SearchByName(ctx context.Context, name string) ([]model.Species, error)
}

type speciesRepo struct {
	db *sql.DB
}

func NewSpeciesRepository(db *sql.DB) SpeciesRepository {
	return &speciesRepo{db: db}
}
