package dto

type MLPrediction struct {
	SpeciesID  int64   `json:"species_id"`
	Confidence float64 `json:"confidence"`
}

type MLResponse struct {
	Predictions   []MLPrediction `json:"predictions"`
	InferenceTime float64        `json:"inference_time_ms"`
}
