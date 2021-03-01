package v1

import (
	"log"
	"net/http"
)

// Init starts the server
func Init() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/book", GetBook)
	mux.HandleFunc("/v1/filters", GetFilters)
	mux.HandleFunc("/v1/ranked_books", GetRankedBooks)

	log.Println("starting server")
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		return err
	}

	return nil
}
