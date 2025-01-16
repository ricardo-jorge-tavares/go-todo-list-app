package middlewares

import (
	"encoding/json"
	"log"
	"net/http"

	"local.com/todo-list-app/internal/config"
	"local.com/todo-list-app/internal/helpers"
)

func AuthMiddleware(next http.Handler) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		contentType := r.Header.Get("Content-Type")

		userName := config.ValidUsers[r.PathValue("userId")]
		if userName == "" {

			if contentType == "application/json" {

				var returnData = struct {
					Error string `json:"error"`
				}{Error: "Invalid user!"}

				w.Header().Set("Content-Type", contentType)
				w.WriteHeader(http.StatusPreconditionFailed)
				json.NewEncoder(w).Encode(returnData)

			} else {

				t, _ := helpers.ParseView("web/views/theme/error.html")
				err := t.Execute(w, "You are trying to access an invalid user")
				if err != nil {
					log.Fatal(err)
				}

			}

			return

		}

		// fmt.Println("Before %s", r.URL.String())
		w.Header().Set("Content-Type", contentType)
		next.ServeHTTP(w, r)
		// fmt.Println("After %s", r.URL.String())

	})

}
