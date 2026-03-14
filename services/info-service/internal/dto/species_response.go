package dto

type CategoryShort struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type SpeciesResponse struct {
	ID             int64    `json:"id" db:"id"`
	ScientificName string   `json:"scientific_name" db:"scientific_name"`
	CommonName     string   `json:"common_name" db:"common_name"`
	Description    string   `json:"description" db:"description"`
	Edibility      int      `json:"edibility" db:"edibility"`
	ToxicityLevel  int      `json:"toxicity_level" db:"toxicity_level"`
	Images         []string `json:"images" db:"images"`

	Category CategoryShort `json:"category" db:"category"`
}

type PaginatedSpeciesResponse struct {
	Data []SpeciesResponse `json:"data"`
	Meta Meta              `json:"meta"`
}

type SpeciesFilter struct {
	Name        *string
	CategoryID  *int64
	Edibility   *int
	ToxicityMax *int

	Sort   string
	Order  string
	Page   int
	Limit  int
	Cursor *int64
}
