package main

import "context"

type MLClient interface {
	Predict(ctx context.Context, image []byte) (*MLResponse, error)
}
