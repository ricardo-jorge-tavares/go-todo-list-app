package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"local.com/todo-list-app/internal/cache"
	"local.com/todo-list-app/internal/helpers"
	"local.com/todo-list-app/internal/services"
	"local.com/todo-list-app/internal/sqldb"
	"local.com/todo-list-app/internal/types"
	"local.com/todo-list-app/web"
)

func main() {

	// Load environment variables from .env
	err := godotenv.Load()
	helpers.CheckError(err)

	// Connect to the database and initialize repositories.
	db := sqldb.ConnectDB()
	defer db.Close()

	// Initialize repositories.
	sqldbTodoRepository := sqldb.NewToDoRepository(db)

	// Create a new Cache instance.
	todoListCache := cache.New[string, types.TodoListItemType]()

	todoService := services.TodoServiceInit(todoListCache, sqldbTodoRepository)

	// Create a new ServeMux instance.
	router := http.NewServeMux()

	// Serve static files.
	fs := http.FileServer(http.Dir("web/assets"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	// Index route.
	router.HandleFunc("GET /{$}", web.IndexRoute)

	// Web (UI) routes.
	appRouter := web.NewAppController(todoService).RegisterRoutes()
	router.Handle("/app/", http.StripPrefix("/app", appRouter))

	fmt.Println("Server started at http://localhost:8080")
	serverError := http.ListenAndServe(":8080", router)
	helpers.CheckError(serverError)

}
