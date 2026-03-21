package dto

type MLPrediction struct {
	ClassID    int64   `json:"species_id"`
	Confidence float64 `json:"probability"`
}

type MLResponse struct {
	Predictions   []MLPrediction `json:"predictions"`
	InferenceTime float64        `json:"inference_time_ms"`
}
