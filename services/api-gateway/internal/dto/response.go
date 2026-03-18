package dto

type PredictionResult struct {
	ID             int64    `json:"id"`
	ScientificName string   `json:"scientific_name"`
	CommonName     string   `json:"common_name"`
	Confidence     float64  `json:"confidence"`
	Images         []string `json:"images"`
}

type PredictResponse struct {
	Data []PredictionResult `json:"data"`
}
