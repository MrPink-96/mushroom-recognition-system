package dto

type CategoryResponse struct {
	ID          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}

type PaginatedCategoryResponse struct {
	Data []CategoryResponse `json:"data"`
	Meta Meta               `json:"meta"`
}
