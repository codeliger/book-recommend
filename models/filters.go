package models

// Book models a basic book
type Filters struct {
	Genres       []string `json:"genres"`
	Authors      []string `json:"authors"`
	MinPageCount int64    `json:"min_page_count"`
	MaxPageCount int64    `json:"max_page_count"`
}
