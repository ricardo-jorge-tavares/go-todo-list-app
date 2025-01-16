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
	r.HandleFunc("POST /login/{$}", appLoginRoute)
	r.HandleFunc("GET /{$}", appIndexRoute)
	r.HandleFunc("GET /{userId}/", c.appListRoute)
	r.HandleFunc("POST /{userId}/{$}", c.appNewTodoRoute)
	return r

}

func appLoginRoute(w http.ResponseWriter, r *http.Request) {

	fmt.Println("/app/login route served")

	r.ParseForm()

	userId := r.FormValue("formUser")
	if userId == "" {
		t, _ := helpers.ParseView("web/views/theme/error.html")
		err := t.Execute(w, "formUser not supplied when trying to login")
		helpers.CheckError(err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/app/%s/", userId), http.StatusFound)
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
		IsCompleted bool
		Rank        int
		CreatedAt   time.Time
	}

	viewData := struct {
		UserId   string
		TodoList []todoListType
	}{
		UserId: userId,
	}

	list := c.todoService.GetUserTodoList(userId)
	for _, item := range list {
		viewData.TodoList = append(viewData.TodoList, todoListType{item.Id, item.Description, item.IsCompleted, item.Rank, item.CreatedAt})
	}

	t, _ := helpers.ParseView("web/views/app/list.html")
	err := t.Execute(w, viewData)
	helpers.CheckError(err)

}

func (c *appController) appNewTodoRoute(w http.ResponseWriter, r *http.Request) {

	fmt.Println("/app/{userId}/new route served")

	userId := r.PathValue("userId")
	fmt.Println(userId)

	r.ParseForm()

	description := r.FormValue("formMessage")
	if description == "" {
		t, _ := helpers.ParseView("web/views/theme/error.html")
		err := t.Execute(w, "formMessage not supplied when trying to insert new Todo item")
		helpers.CheckError(err)
		return
	}

	c.todoService.AddTodoItem(userId, description)

	http.Redirect(w, r, fmt.Sprintf("/app/%s/", userId), http.StatusFound)

}
