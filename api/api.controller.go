package api

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	r.HandleFunc("POST /{userId}/todo/{todoId}/description/{$}", c.apiUpdateDescription)
	r.HandleFunc("POST /{userId}/todo/{todoId}/rank/{$}", c.apiUpdateRank)

	return r
}

func (c *apiController) apiUpdateDescription(w http.ResponseWriter, r *http.Request) {

	fmt.Println("/app/{userId}/todo/{todoId}/description route served")

	userId := r.PathValue("userId")
	todoId := r.PathValue("todoId")
	fmt.Println(userId, todoId)

	var inputData struct {
		Description string `json:"description"`
	}
	json.NewDecoder(r.Body).Decode(&inputData)

	c.todoService.UpdateTodoItemDescription(userId, todoId, inputData.Description)

	var returnData = struct {
		Id string `json:"id"`
	}{
		Id: todoId,
	}
	json.NewEncoder(w).Encode(returnData)

}

func (c *apiController) apiUpdateRank(w http.ResponseWriter, r *http.Request) {

	fmt.Println("/api/{userId}/todo/{todoId}/rank route served")

	userId := r.PathValue("userId")
	todoId := r.PathValue("todoId")
	fmt.Println(userId, todoId)

	var inputData struct {
		Rank int `json:"rank"`
	}
	json.NewDecoder(r.Body).Decode(&inputData)

	var returnData = struct {
		Id           string `json:"id"`
		ErrorMessage string `json:"error"`
	}{
		Id: todoId,
	}

	c.todoService.UpdateTodoItemRank(userId, todoId, inputData.Rank)

	json.NewEncoder(w).Encode(returnData)

}
