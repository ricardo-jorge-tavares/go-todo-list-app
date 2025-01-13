package web

import (
	"fmt"
	"net/http"

	"local.com/todo-list-app/internal/helpers"
	"local.com/todo-list-app/internal/services"
)

type appController struct {
	// todoService *services.TodoService
	todoService services.TodoServiceInterface
}

func NewAppController(s services.TodoServiceInterface) *appController {
	return &appController{
		todoService: s,
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

	var viewData struct {
		TodoList []struct {
			id          string
			Description string
		}
	}

	list := c.todoService.GetTodoList(userId)
	for _, todo := range list {
		viewData.TodoList = append(viewData.TodoList, struct {
			id          string
			Description string
		}{todo.Id, todo.Description})
	}

	t, _ := helpers.ParseView("web/views/app/list.html")
	err := t.Execute(w, viewData)
	helpers.CheckError(err)

}
