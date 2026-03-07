package dto

type Meta struct {
	Page       int    `json:"page,omitempty"`
	Limit      int    `json:"limit"`
	Total      int    `json:"total,omitempty"`
	Pages      int    `json:"pages,omitempty"`
	NextCursor *int64 `json:"next_cursor,omitempty"`
}
