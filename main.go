package main

import (
	"fmt"
	"net/http"

	helpers "local.com/todo-list-app/internal/helpers"
	web "local.com/todo-list-app/web"
)

func main() {

	router := http.NewServeMux()

	// Serve static files.
	fs := http.FileServer(http.Dir("web/assets"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	// Index route.
	router.HandleFunc("GET /{$}", web.IndexRoute)

	// UI routes.
	appRouter := web.NewAppHandler().RegisterRoutes()
	router.Handle("/app/", http.StripPrefix("/app", appRouter))

	fmt.Println("Server started at http://localhost:8080")
	serverError := http.ListenAndServe(":8080", router)
	helpers.CheckError(serverError)

}
