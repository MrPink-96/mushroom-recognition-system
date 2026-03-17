package main

import "context"

type InfoClient interface {
	GetByIDs(ctx context.Context, ids []int64) ([]Species, error)
}
