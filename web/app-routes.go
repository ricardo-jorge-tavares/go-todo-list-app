package web

import (
	"net/http"

	helpers "local.com/todo-list-app/internal/helpers"
)

type AppHandler struct{}

func NewAppHandler() *AppHandler {
	return &AppHandler{}
}

func (h *AppHandler) RegisterRoutes() *http.ServeMux {

	r := http.NewServeMux()
	r.HandleFunc("GET /{$}", appListRoute)
	// r.HandleFunc("POST /", appListRoute)
	// r.HandleFunc("GET /{id}/", appListRoute)
	// r.HandleFunc("DELETE /{id}/", appListRoute)

	return r
}

func appListRoute(writer http.ResponseWriter, request *http.Request) {

	t, _ := helpers.ParseView("web/views/app/list.html")
	err := t.Execute(writer, nil)
	helpers.CheckError(err)

	// fmt.Printf("Key: %d, Value: %s\n", writer, request.URL)
	// // Create a template using the html
	// tmpl, err := template.ParseFiles("web/views/app/list.html")
	// helpers.CheckError(err)

	// err = tmpl.Execute(writer, nil)
	// helpers.CheckError(err)
}
