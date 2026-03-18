package client

import (
	"api-gateway/internal/dto"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	appErr "api-gateway/internal/errors"
)

type InfoClient interface {
	GetByIDs(ctx context.Context, ids []int64) ([]dto.PredictionResult, error)
	Health(ctx context.Context) error
}

type infoClient struct {
	baseURL string
	client  *http.Client
}

func NewInfoClient(baseURL string, client *http.Client) InfoClient {
	return &infoClient{baseURL:baseURL, client	: client}
}

func (c *infoClient) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/health", nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return appErr.ErrInfoUnavailable
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("Health check failed: status=%d, body=%s", resp.StatusCode, body)

		if resp.StatusCode == http.StatusServiceUnavailable {
			return appErr.ErrInfoUnavailable
		}

		return fmt.Errorf("unexpected health check response: %d", resp.StatusCode)
	}
	return nil
}

func (c *infoClient) GetByIDs(ctx context.Context, ids []int64) ([]dto.PredictionResult, error)


	return
}
