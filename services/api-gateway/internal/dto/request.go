package dto

import "mime/multipart"

type PredictRequest struct {
	Image    *multipart.FileHeader `form:"image"`
	ImageURL *string               `json:"image_url"`
}
