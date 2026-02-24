package model

type Species struct {
	ID                int64  `json:"id"`
	ScientificName    string `json:"scientific_name"`
	CommonName        string `json:"common_name"`
	Description       string `json:"description"`
	Edibility         int    `json:"edibility"`
	ToxicityLevel     int    `json:"toxicity_level"`
	ReferenceImageURL string `json:"reference_image_url"`
	CategoryID        int64  `json:"category_id"`
}
