package models

// Book models a basic book
type Book struct {
	ID              int64    `json:"id"`
	Title           string   `json:"title"`
	Authors         []string `json:"authors"`
	Genres          []string `json:"genres"`
	Pages           int64    `json:"pages"`
	PublicationYear int64    `json:"publication_year"`
	Rating          int64    `json:"rating"`
}
