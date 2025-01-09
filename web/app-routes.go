package web

import (
	"fmt"
	"net/http"
	"time"

	"local.com/todo-list-app/internal/cache"
	"local.com/todo-list-app/internal/helpers"
	"local.com/todo-list-app/internal/sqldb"
	"local.com/todo-list-app/internal/types"
)

type appController struct {
	cache     *cache.Cache[string, types.TodoListItemType]
	sqldbTodo *sqldb.SqlToDoRepository
}

func NewAppController(c *cache.Cache[string, types.TodoListItemType], s *sqldb.SqlToDoRepository) *appController {
	return &appController{
		cache:     c,
		sqldbTodo: s,
	}
}

func (c *appController) RegisterRoutes() *http.ServeMux {

	r := http.NewServeMux()
	r.HandleFunc("GET /{$}", appIndexRoute)
	r.HandleFunc("GET /{userId}", c.appListRoute)
	// r.HandleFunc("GET /{id}/", appListRoute)
	// r.HandleFunc("DELETE /{id}/", appListRoute)

	return r
}

func appIndexRoute(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/app/user123", http.StatusFound)
}

func (c *appController) appListRoute(w http.ResponseWriter, r *http.Request) {

	userId := r.PathValue("userId")
	fmt.Println(userId)

	c.cache.Set(userId, types.TodoListItemType{Description: "Go to the gym", CreatedAt: time.Now(), IsComplete: false})

	for k, v := range c.cache.List() {
		fmt.Printf("Key: %s, Value: %s | %s | %v\n", k, v.Description, v.CreatedAt, v.IsComplete)
	}

	sqlId := c.sqldbTodo.Insert("Go to the gym " + userId)
	fmt.Println("Inserted record with ID:", sqlId)

	row, _ := c.sqldbTodo.FindAll()

	for _, todo := range row {
		fmt.Println("From DB", todo.Id, todo.Description)
	}

	t, _ := helpers.ParseView("web/views/app/list.html")
	err := t.Execute(w, nil)
	helpers.CheckError(err)

}
