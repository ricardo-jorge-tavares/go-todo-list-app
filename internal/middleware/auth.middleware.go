package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"local.com/todo-list-app/internal/config"
	"local.com/todo-list-app/internal/helpers"
	"local.com/todo-list-app/internal/models"
)

func AuthMiddleware(next http.Handler) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		contentType := r.Header.Get("Content-Type")

		userId := r.PathValue("userId")
		userName := config.ValidUsers[userId]

		if userName == "" {

			if contentType == "application/json" {

				var returnData = struct {
					Error string `json:"error"`
				}{Error: "Invalid user!"}

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

		userContext := models.UserContext{Id: userId, Name: userName}

		rcopy := r.WithContext(models.SetContextUser(r.Context(), &userContext))
		next.ServeHTTP(w, rcopy)

	})

}
