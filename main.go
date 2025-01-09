package main

import (
	"fmt"
	"net/http"

	"local.com/todo-list-app/internal/cache"
	"local.com/todo-list-app/internal/helpers"
	"local.com/todo-list-app/internal/types"
	"local.com/todo-list-app/web"
)

func main() {

	// Create a new Cache instance
	todoListCache := cache.New[string, types.TodoListItemType]()

	router := http.NewServeMux()

	// Serve static files.
	fs := http.FileServer(http.Dir("web/assets"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	// Index route.
	router.HandleFunc("GET /{$}", web.IndexRoute)

	// UI routes.
	appRouter := web.NewAppController(todoListCache).RegisterRoutes()
	router.Handle("/app/", http.StripPrefix("/app", appRouter))

	fmt.Println("Server started at http://localhost:8080")
	serverError := http.ListenAndServe(":8080", router)
	helpers.CheckError(serverError)

}
