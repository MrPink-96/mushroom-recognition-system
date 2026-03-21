package service

import (
	"api-gateway/internal/client"
	"api-gateway/internal/dto"
	"context"
	"net/http"
)

const topK = 5

type PredictService struct {
	mlClient   client.MLClient
	infoClient client.InfoClient
}

func NewPredictService(mlClient client.MLClient, infoClient client.InfoClient) *PredictService {
	return &PredictService{mlClient: mlClient, infoClient: infoClient}
}

func (s *PredictService) Predict(ctx context.Context, req *http.Request) (*dto.PredictResponse, error) {
	mlResp, err := s.mlClient.Predict(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(mlResp.Predictions) == 0 {
		return &dto.PredictResponse{Data: []dto.PredictionResult{}}, nil
	}

	limit := topK
	if len(mlResp.Predictions) < limit {
		limit = len(mlResp.Predictions)
	}

	rawPreds := mlResp.Predictions[:limit]
	seen := make(map[int64]struct{})
	preds := make([]dto.MLPrediction, 0, len(rawPreds))

	for _, p := range rawPreds {
		if p.ClassID <= 0 || p.Confidence <= 0 {
			continue
		}

		if _, exists := seen[p.ClassID]; exists {
			continue
		}

		seen[p.ClassID] = struct{}{}
		preds = append(preds, p)
	}

	if len(preds) == 0 {
		return &dto.PredictResponse{Data: []dto.PredictionResult{}}, nil
	}

	ids := make([]int64, 0, len(preds))
	for _, p := range preds {
		ids = append(ids, p.ClassID)
	}

	infResp, err := s.infoClient.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	idToSpecies := make(map[int64]dto.PredictionResult, len(infResp))
	for _, sp := range infResp {
		idToSpecies[sp.ID] = sp
	}

	result := make([]dto.PredictionResult, 0, len(preds))

	for _, p := range preds {
		sp, ok := idToSpecies[p.ClassID]
		if !ok {
			continue
		}

		sp.Confidence = p.Confidence
		result = append(result, sp)
	}

	return &dto.PredictResponse{Data: result}, nil
}
