package helpers

import (
	"log"
	"net/http"
)

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

// func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusNotFound)
// 	w.Write([]byte("404 Not Found"))
// }

func CheckError(err error) {
	// Handle errors
	if err != nil {
		log.Fatal(err)
	}
}
