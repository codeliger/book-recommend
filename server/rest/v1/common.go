package v1

import (
	"fmt"
	"log"
	"net/http"
)

// WriteError writes an error to a response writer
func WriteError(w http.ResponseWriter, status int, friendlyError string, optionalErr error) {
	log.Println(friendlyError, optionalErr)

	w.WriteHeader(status)

	_, err := w.Write([]byte(friendlyError))

	if err != nil {
		fmt.Println("failed to write body ", err)
	}
}
