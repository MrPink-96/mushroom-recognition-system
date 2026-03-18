package dto

type MLPrediction struct {
	ClassID    int64   `json:"class_id"`
	Confidence float64 `json:"confidence"`
}

type MLResponse struct {
	Predictions   []MLPrediction `json:"predictions"`
	InferenceTime float64        `json:"inference_time_ms"`
}
