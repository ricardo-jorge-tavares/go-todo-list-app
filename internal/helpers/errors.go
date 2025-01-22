package helpers

import (
	"errors"
	"log"
	"net/http"

	"local.com/todo-list-app/internal/models"
)

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"error": "Internal Server Error"}`))
}

func PreconditionFailedHandler(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusPreconditionFailed)
	w.Write([]byte(`{"error": "` + message + `"}`))
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

func HandlerError(w http.ResponseWriter, err error) {

	// Still not finished! Should have logic to support a generic error!
	var e models.ErrorModel
	if errors.As(err, &e) {

		switch e.GetType() {
		case models.ErrorBadRequest:
			w.WriteHeader(http.StatusPreconditionFailed)
		case models.ErrorInternalFailure:
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write([]byte(`{
			"errorCode": "` + e.GetCode() + `",
			"errorMessage": "` + e.GetMessage() + `"
		}`))

	}
}
