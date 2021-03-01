package db

import (
	"encoding/json"
	"io/ioutil"
	"sort"

	"github.com/codeliger/book-recommend/models"
)

var books []models.Book
var filters models.Filters

// Init reads the books into memory
func Init() error {

	data, err := ioutil.ReadFile("books.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &books)

	if err != nil {
		return err
	}

	filters = InitFilters()

	return nil
}

// SelectBook selects a book by id
func SelectBook(id int64) *models.Book {
	for _, book := range books {
		if book.ID == id {
			return &book
		}
	}

	return nil
}

// SelectBooks selects all books
func SelectBooks() []models.Book {
	return books
}

// InitFilters creates a list of genres and authors; so the index can be reliably used
func InitFilters() models.Filters {
	minPageCount := 0
	maxPageCount := 0

	books := SelectBooks()

	uniqueAuthors := []string{}
	uniqueGenres := []string{}

	for _, book := range books {

		for _, bookAuthor := range book.Authors {
			inUniqueAuthors := false
			for _, uniqueAuthor := range uniqueAuthors {
				if bookAuthor == uniqueAuthor {
					inUniqueAuthors = true
				}
			}
			if !inUniqueAuthors {
				uniqueAuthors = append(uniqueAuthors, bookAuthor)
			}
		}

		for _, bookGenre := range book.Genres {
			inUniqueGenres := false
			for _, uniqueGenre := range uniqueGenres {
				if bookGenre == uniqueGenre {
					inUniqueGenres = true
				}
			}
			if !inUniqueGenres {
				uniqueGenres = append(uniqueGenres, bookGenre)
			}
		}

		if book.Pages < int64(minPageCount) {
			minPageCount = int(book.Pages)
		}

		if book.Pages > int64(maxPageCount) {
			maxPageCount = int(book.Pages)
		}
	}

	sort.Strings(uniqueGenres)
	sort.Strings(uniqueAuthors)

	tempFilters := models.Filters{
		Genres:       uniqueGenres,
		Authors:      uniqueAuthors,
		MinPageCount: int64(minPageCount),
		MaxPageCount: int64(maxPageCount),
	}

	return tempFilters
}

// GetFilters returns the initialized filters used to reference authors and genres by order
func GetFilters() models.Filters {
	return filters
}
