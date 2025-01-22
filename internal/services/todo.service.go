package services

import (
	"fmt"
	"sort"
	"time"

	"local.com/todo-list-app/internal/cache"
	"local.com/todo-list-app/internal/models"
	"local.com/todo-list-app/internal/sqldb"
)

// Structs.
type TodoService struct {
	cache         *cache.Cache[string, models.CacheUserModel]
	sqlDbTodoRepo sqldb.ToDoRepositoryInterface
}

type GetUserTodoListResponse struct {
	Id          string    `json:"id"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"isCompleted"`
	Rank        int       `json:"rank"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Interfaces.
type TodoServiceInterface interface {
	GetUserTodoList(userId string) []GetUserTodoListResponse
	AddTodoItem(userId string, description string) (id string, err error)
	UpdateTodoItemDescription(userId string, todoId string, description string) (id string)
	UpdateTodoItemRank(userId string, todoId string, rank int) (id string)
	UpdateTodoItemIsCompleted(userId string, todoId string) (id string)
	DeleteTodoItem(userId string, todoId string) (id string)
}

// Functions.
func TodoServiceInit(c *cache.Cache[string, models.CacheUserModel], sqldbTodo sqldb.ToDoRepositoryInterface) *TodoService {
	return &TodoService{
		cache:         c,
		sqlDbTodoRepo: sqldbTodo,
	}
}

func (s *TodoService) GetUserTodoList(userId string) (r []GetUserTodoListResponse) {

	user, found := s.cache.Get(userId)

	// Check if the user is found in cache and if is still valid.
	if found && user.ExpiresAt.After(time.Now()) {
		fmt.Println("User found and valid in cache. Returning it!")
		for k, v := range user.TodoList.List() {
			r = append(r, GetUserTodoListResponse{Id: k, Description: v.Description, IsCompleted: v.IsCompleted, Rank: v.Rank, CreatedAt: v.CreatedAt})
		}

		sort.Slice(r, func(i, j int) bool {
			return r[i].Rank < r[j].Rank
		})

		return r
	}

	fmt.Println("Fetching info From DB")
	rows, _ := s.sqlDbTodoRepo.FindUserItems(userId)

	user = s.cacheSetUser(userId)

	for _, item := range rows {

		user.TodoList.Set(item.Id, models.CacheTodoItemModel{
			Description: item.Description,
			IsCompleted: item.IsCompleted,
			Rank:        item.Rank,
			CreatedAt:   item.CreatedAt,
		})

		r = append(r, GetUserTodoListResponse{Id: item.Id, Description: item.Description, IsCompleted: item.IsCompleted, Rank: item.Rank, CreatedAt: item.CreatedAt})

	}

	return r

}

func (s *TodoService) AddTodoItem(userId string, description string) (id string, err error) {

	if userId == "" {
		return "", models.NewError(models.ErrorBadRequest, "S.TD.001", "User Id is required")
	}
	if description == "" {
		return "", models.NewError(models.ErrorBadRequest, "S.TD.002", "description is required")
	}

	sqlTodoId := s.sqlDbTodoRepo.InsertItem(userId, description)
	// fmt.Println("Inserted record with ID:", sqlTodoId)

	s.cacheInvalidateUser(userId)

	return sqlTodoId, nil

}

func (s *TodoService) UpdateTodoItemDescription(userId string, todoId string, description string) (id string) {

	s.sqlDbTodoRepo.UpdateItemDescription(todoId, description)

	s.cacheInvalidateUser(userId)

	return todoId

}

func (s *TodoService) UpdateTodoItemRank(userId string, todoId string, rank int) (id string) {

	s.sqlDbTodoRepo.UpdateItemRank(userId, todoId, rank)

	s.cacheInvalidateUser(userId)

	return todoId

}

func (s *TodoService) UpdateTodoItemIsCompleted(userId string, todoId string) (id string) {

	s.sqlDbTodoRepo.UpdateItemIsCompleted(todoId)

	s.cacheInvalidateUser(userId)

	return todoId

}

func (s *TodoService) DeleteTodoItem(userId string, todoId string) (id string) {

	s.sqlDbTodoRepo.DeleteItem(todoId)

	s.cacheInvalidateUser(userId)

	return todoId

}

func (s *TodoService) cacheSetUser(userId string) models.CacheUserModel {

	return s.cache.Set(userId, models.CacheUserModel{
		TodoList:  cache.New[string, models.CacheTodoItemModel](),
		ExpiresAt: time.Now().Add(30 * time.Second),
	})

}

func (s *TodoService) cacheInvalidateUser(userId string) {

	s.cache.Delete(userId)

}
