package middlewares

import (
	"fmt"
	"net/http"
)

func ApiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("Before %s\n", r.URL.String())

		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)

		// fmt.Printf("After %s", r.URL.String())
	})
}
