package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"local.com/todo-list-app/internal/helpers"
	"local.com/todo-list-app/internal/middleware"
	"local.com/todo-list-app/internal/services"
)

type apiController struct {
	todoService services.TodoServiceInterface
}

func NewApiController(s services.TodoServiceInterface) *apiController {
	return &apiController{
		todoService: s,
	}
}

func (c *apiController) RegisterRoutes() *http.ServeMux {

	r := http.NewServeMux()

	var handlers = map[string]http.HandlerFunc{
		"GET /{userId}/": c.apiListUserTodos,
		"POST /{userId}/todo/{todoId}/description/{$}": c.apiUpdateTodoDescription,
		"POST /{userId}/todo/{todoId}/rank/{$}":        c.apiUpdateTodoRank,
		"POST /{userId}/todo/{todoId}/completed/{$}":   c.apiUpdateTodoIsCompleted,
		"DELETE /{userId}/todo/{todoId}/{$}":           c.apiDeleteTodo,
	}

	// Apply middewares at a router level (so that middleware can access path params)
	for pattern, handler := range handlers {
		r.HandleFunc(pattern, middleware.AuthMiddleware(handler))
	}

	return r
}

func (c *apiController) apiListUserTodos(w http.ResponseWriter, r *http.Request) {

	fmt.Println("/api/{userId} route served")

	userId := r.PathValue("userId")

	list := c.todoService.GetUserTodoList(userId)
	returnData, err := json.Marshal(list)
	if err != nil {
		helpers.InternalServerErrorHandler(w, r)
		return
	}

	w.Write(returnData)

}

func (c *apiController) apiUpdateTodoDescription(w http.ResponseWriter, r *http.Request) {

	fmt.Println("/app/{userId}/todo/{todoId}/description route served")

	userId := r.PathValue("userId")
	todoId := r.PathValue("todoId")

	var inputData struct {
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		helpers.InternalServerErrorHandler(w, r)
		return
	}

	c.todoService.UpdateTodoItemDescription(userId, todoId, inputData.Description)

	var returnData = struct {
		Id string `json:"id"`
	}{
		Id: todoId,
	}

	json.NewEncoder(w).Encode(returnData)

}

func (c *apiController) apiUpdateTodoRank(w http.ResponseWriter, r *http.Request) {

	fmt.Println("/api/{userId}/todo/{todoId}/rank route served")

	userId := r.PathValue("userId")
	todoId := r.PathValue("todoId")

	var inputData struct {
		Rank int `json:"rank"`
	}

	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		helpers.InternalServerErrorHandler(w, r)
		return
	}

	var returnData = struct {
		Id           string `json:"id"`
		ErrorMessage string `json:"error"`
	}{
		Id: todoId,
	}

	c.todoService.UpdateTodoItemRank(userId, todoId, inputData.Rank)

	json.NewEncoder(w).Encode(returnData)

}

func (c *apiController) apiUpdateTodoIsCompleted(w http.ResponseWriter, r *http.Request) {

	userId := r.PathValue("userId")
	todoId := r.PathValue("todoId")

	var returnData = struct {
		Id           string `json:"id"`
		ErrorMessage string `json:"error"`
	}{
		Id: todoId,
	}

	c.todoService.UpdateTodoItemIsCompleted(userId, todoId)

	json.NewEncoder(w).Encode(returnData)

}

func (c *apiController) apiDeleteTodo(w http.ResponseWriter, r *http.Request) {

	userId := r.PathValue("userId")
	todoId := r.PathValue("todoId")

	var returnData = struct {
		Id           string `json:"id"`
		ErrorMessage string `json:"error"`
	}{
		Id: todoId,
	}

	c.todoService.DeleteTodoItem(userId, todoId)

	json.NewEncoder(w).Encode(returnData)

}
