package dto

type CategoryShort struct {
	ID   int64  `json:"id" db:"category.id"`
	Name string `json:"name" db:"category.name"`
}

type SpeciesResponse struct {
	ID                int64         `json:"id" db:"id"`
	ScientificName    string        `json:"scientific_name" db:"scientific_name"`
	CommonName        string        `json:"common_name" db:"common_name"`
	Description       string        `json:"description" db:"description"`
	Edibility         int           `json:"edibility" db:"edibility"`
	ToxicityLevel     int           `json:"toxicity_level" db:"toxicity_level"`
	ReferenceImageURL string        `json:"reference_image_url" db:"reference_image_url"`
	Category          CategoryShort `json:"category"`
}

type Meta struct {
	Page       int    `json:"page,omitempty"`
	Limit      int    `json:"limit"`
	Total      int    `json:"total,omitempty"`
	Pages      int    `json:"pages,omitempty"`
	NextCursor *int64 `json:"next_cursor,omitempty"`
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
