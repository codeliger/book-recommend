package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/codeliger/book-recommend/db"
	"github.com/codeliger/book-recommend/models"
)

// GetBook handles get requests for books
func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	strid := r.URL.Query().Get("id")

	if strid == "" {
		books := db.SelectBooks()
		data, err := json.Marshal(books)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "failed to return decoded list of books", err)
			return
		}
		w.Header().Add("Content-Type", "text/json")
		w.Write(data)
		return
	}

	id, err := strconv.Atoi(strid)

	if err != nil {
		WriteError(w, http.StatusNotFound, "id is not valid integer", err)
		return
	}

	int64ID := int64(id)

	book := db.SelectBook(int64ID)

	if book == nil {
		WriteError(w, http.StatusNotFound, "id does not exist", nil)
		return
	}

	data, err := json.Marshal(book)

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to return decoded book", err)
		return
	}
	w.Header().Add("Content-Type", "text/json")
	w.Write(data)
}

// GetFilters returns the filters to be used on the book recommendation website
func GetFilters(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	filters := db.GetFilters()

	data, err := json.Marshal(filters)

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to enumerate filters", err)
		return
	}

	w.Header().Add("Content-Type", "text/json")
	_, err = w.Write(data)

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to marshal data; encoding problem", err)
		return
	}
}

func GetRankedBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	query := r.URL.Query()

	var authorIDs []int64
	var genreIDs []int64
	var minPageCount int64
	var maxPageCount int64

	authorString := query.Get("authors")

	if authorString != "" {
		authorStringIDs := strings.Split(authorString, ":")
		for _, stringID := range authorStringIDs {
			authorID, err := strconv.ParseInt(stringID, 10, 64)
			if err != nil {
				WriteError(w, http.StatusInternalServerError, "failed to parse stringID", err)
				return
			}
			authorIDs = append(authorIDs, authorID)
		}
	}
	genresString := query.Get("genres")

	if genresString != "" {
		genreStringIDs := strings.Split(genresString, ":")
		for _, stringID := range genreStringIDs {
			genreID, err := strconv.ParseInt(stringID, 10, 64)
			if err != nil {
				WriteError(w, http.StatusInternalServerError, "failed to parse stringID", err)
				return
			}
			genreIDs = append(genreIDs, genreID)
		}
	}

	minPageCountString := query.Get("min_page_count")

	var err error

	if minPageCountString != "" {
		minPageCount, err = strconv.ParseInt(minPageCountString, 10, 64)

		if err != nil {
			WriteError(w, http.StatusInternalServerError, "failed to parse minPageCount", err)
			return
		}
	}
	maxPageCountString := query.Get("max_page_count")

	if minPageCountString != "" {
		maxPageCount, err = strconv.ParseInt(maxPageCountString, 10, 64)

		if err != nil {
			WriteError(w, http.StatusInternalServerError, "failed to parse maxPageCount", err)
			return
		}
	}

	books := db.SelectBooks()
	filters := db.GetFilters()

	var submittedAuthorNames []string
	var submittedGenreNames []string

	// because the books store the names and not the ids:

	// convert ids to names from filter lookup table
	authorLookupLength := len(filters.Authors)
	for _, submittedAuthorID := range authorIDs {
		if submittedAuthorID < int64(authorLookupLength) {
			filterAuthorName := filters.Authors[int(submittedAuthorID)]
			submittedAuthorNames = append(submittedAuthorNames, filterAuthorName)
		}
	}

	// convert ids to names from filter lookup table
	genreLookupLength := len(filters.Genres)
	for _, submittedGenreID := range genreIDs {
		if submittedGenreID < int64(genreLookupLength) {
			filterGenreName := filters.Genres[int(submittedGenreID)]
			submittedGenreNames = append(submittedGenreNames, filterGenreName)
		}
	}

	authorFilteredBooks := make([]models.Book, 0)
	genreFilteredBooks := make([]models.Book, 0)

	if len(submittedAuthorNames) > 0 {

		for _, book := range books {
			for _, bookauthor := range book.Authors {
				authorInFilters := false
				for _, submittedAuthor := range submittedAuthorNames {
					if bookauthor == submittedAuthor {
						authorInFilters = true
					}
				}
				if authorInFilters {
					authorFilteredBooks = append(authorFilteredBooks, book)
				}
			}
		}

	} else {
		authorFilteredBooks = books
	}

	log.Println("authorFilteredBooks", authorFilteredBooks)

	if len(submittedGenreNames) > 0 {

		for _, book := range authorFilteredBooks {
			for _, bookgenre := range book.Genres {
				genreInFilters := false
				for _, submittedGenre := range submittedGenreNames {
					if bookgenre == submittedGenre {
						genreInFilters = true
					}
				}
				if genreInFilters {
					genreFilteredBooks = append(genreFilteredBooks, book)
				}
			}

		}
	} else {
		genreFilteredBooks = authorFilteredBooks
	}

	log.Println("genreFilteredBooks", genreFilteredBooks)

	minPageCountFilteredBooks := make([]models.Book, 0)

	if minPageCountString != "" {
		for _, book := range genreFilteredBooks {
			if book.Pages >= minPageCount {
				minPageCountFilteredBooks = append(minPageCountFilteredBooks, book)
			}
		}
	} else {
		minPageCountFilteredBooks = genreFilteredBooks
	}

	maxPageCountFilteredBooks := make([]models.Book, 0)

	if maxPageCountString != "" {

		for _, book := range minPageCountFilteredBooks {
			if book.Pages <= maxPageCount {
				maxPageCountFilteredBooks = append(maxPageCountFilteredBooks, book)
			}
		}

	} else {
		maxPageCountFilteredBooks = minPageCountFilteredBooks
	}

	log.Println("maxPageCountFilteredBooks", maxPageCountFilteredBooks)

	data, err := json.Marshal(maxPageCountFilteredBooks)

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to marshal json", err)
		return
	}

	w.Header().Add("Content-Type", "text/json")
	w.Write(data)
}
