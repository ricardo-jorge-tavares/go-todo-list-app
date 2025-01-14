package web

import (
	"fmt"
	"net/http"
	"time"

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
	r.HandleFunc("GET /{userId}/", c.appListRoute)
	r.HandleFunc("GET /{userId}/new/", c.appNewRoute)
	// r.HandleFunc("DELETE /{id}/", appListRoute)

	return r
}

func appIndexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/app route served")
	http.Redirect(w, r, "app/user123/", http.StatusFound)
}

func (c *appController) appListRoute(w http.ResponseWriter, r *http.Request) {

	fmt.Println("/app/{userId} route served")

	userId := r.PathValue("userId")
	fmt.Println(userId)

	type todoListType struct {
		Id          string
		Description string
		CreatedAt   time.Time
		IsComplete  bool
	}

	var viewData struct {
		TodoList []todoListType
	}

	list := c.todoService.GetUserTodoList(userId)
	for _, item := range list {
		viewData.TodoList = append(viewData.TodoList, todoListType{item.Id, item.Description, item.CreatedAt, item.IsComplete})
	}

	t, _ := helpers.ParseView("web/views/app/list.html")
	err := t.Execute(w, viewData)
	helpers.CheckError(err)

}

func (c *appController) appNewRoute(w http.ResponseWriter, r *http.Request) {

	fmt.Println("/app/{userId}/new route served")

	userId := r.PathValue("userId")
	fmt.Println(userId)

	var viewData struct {
		TodoList []struct {
			id          string
			Description string
		}
	}

	c.todoService.AddTodoItem(userId, "New todo item"+userId)

	// http.Redirect(w, r, "/app/user123", http.StatusFound)

	t, _ := helpers.ParseView("web/views/app/list.html")
	err := t.Execute(w, viewData)
	helpers.CheckError(err)

}
