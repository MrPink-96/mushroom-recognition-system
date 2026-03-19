package client

import (
	"api-gateway/internal/dto"
	appErr "api-gateway/internal/errors"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type InfoClient interface {
	GetByIDs(ctx context.Context, ids []int64) ([]dto.PredictionResult, error)
	Health(ctx context.Context) error
	ProxyGet(ctx context.Context, path string, query url.Values) ([]byte, int, error)
}

type infoClient struct {
	baseURL string
	client  *http.Client
}

func NewInfoClient(baseURL string, client *http.Client) InfoClient {
	return &infoClient{baseURL: baseURL, client: client}
}

func (c *infoClient) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/health", nil)
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

func (c *infoClient) GetByIDs(ctx context.Context, ids []int64) ([]dto.PredictionResult, error) {
	if len(ids) == 0 {
		return []dto.PredictionResult{}, nil
	}

	idsStr := make([]string, 0, len(ids))
	for _, id := range ids {
		idsStr = append(idsStr, strconv.FormatInt(id, 10))
	}

	q := url.Values{}
	q.Set("ids", strings.Join(idsStr, ","))

	u := c.baseURL + "/species/batch?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, appErr.ErrInfoUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Predict failed: status=%d, body=%s", resp.StatusCode, body)

		if resp.StatusCode == http.StatusServiceUnavailable {
			return nil, appErr.ErrInfoUnavailable
		}
		return nil, fmt.Errorf("predict failed with status %d", resp.StatusCode)
	}

	var result []dto.PredictionResult
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *infoClient) ProxyGet(ctx context.Context, path string, query url.Values) ([]byte, int, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, 0, err
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, 0, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, appErr.ErrInfoUnavailable
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	return body, resp.StatusCode, nil
}
