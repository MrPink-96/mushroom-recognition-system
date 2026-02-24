package dto

type SpeciesResponse struct {
	ID                int64  `db:"id" json:"id"`
	ScientificName    string `db:"scientific_name" json:"scientific_name"`
	CommonName        string `db:"common_name" json:"common_name"`
	Description       string `db:"description" json:"description"`
	Edibility         int    `db:"edibility" json:"edibility"`
	ToxicityLevel     int    `db:"toxicity_level" json:"toxicity_level"`
	ReferenceImageURL string `db:"reference_image_url" json:"reference_image_url"`
	CategoryID        int64  `db:"category_id" json:"category_id"`
	CategoryName      string `db:"category_name" json:"category_name"`
}
