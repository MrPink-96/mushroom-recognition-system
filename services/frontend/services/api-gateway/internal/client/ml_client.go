package client

import (
	"api-gateway/internal/dto"
	appErr "api-gateway/internal/errors"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/health", nil)
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
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "image.jpg")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	_, err = part.Write(image)
	if err != nil {
		return nil, fmt.Errorf("failed to write image: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/predict", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, appErr.ErrMLUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Predict failed: status=%d, body=%s", resp.StatusCode, body)

		if resp.StatusCode == http.StatusServiceUnavailable {
			return nil, appErr.ErrMLUnavailable
		}
		return nil, fmt.Errorf("predict failed with status %d", resp.StatusCode)
	}

	var result dto.MLResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
