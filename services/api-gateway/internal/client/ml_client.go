package client

import (
	"api-gateway/internal/dto"
	appErr "api-gateway/internal/errors"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

type MLClient interface {
	Predict(ctx context.Context, image []byte) (*dto.MLResponse, error)
	Health(ctx context.Context) error
}

type mlClient struct {
	baseURL string
	client  *http.Client
}

func NewMLClient(baseURL string, client *http.Client) MLClient {
	return &mlClient{baseURL: baseURL, client: client}
}

func (c *mlClient) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/health", nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return appErr.ErrMLUnavailable
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("Health check failed: status=%d, body=%s", resp.StatusCode, body)

		if resp.StatusCode == http.StatusServiceUnavailable {
			return appErr.ErrMLUnavailable
		}

		return fmt.Errorf("unexpected health check response: %d", resp.StatusCode)
	}
	return nil

}

func (c *mlClient) Predict(ctx context.Context, image []byte) (*dto.MLResponse, error) {

}
