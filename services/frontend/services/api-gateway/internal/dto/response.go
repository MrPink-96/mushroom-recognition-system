package dto

type PredictionResult struct {
	ID             int64    `json:"id"`
	ScientificName string   `json:"scientific_name"`
	CommonName     string   `json:"common_name"`
	Description    string   `json:"description"`
	Edibility      int      `json:"edibility"`
	ToxicityLevel  int      `json:"toxicity_level"`
	Images         []string `json:"images"`

	Category CategoryShort `json:"category"`

	Confidence float64 `json:"confidence"`
}

type PredictResponse struct {
	Data []PredictionResult `json:"data"`
}
