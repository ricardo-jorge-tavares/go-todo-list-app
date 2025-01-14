package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"local.com/todo-list-app/internal/cache"
	"local.com/todo-list-app/internal/models"
	"local.com/todo-list-app/internal/sqldb"
)

// Structs.
type TodoService struct {
	cache         *cache.Cache[string, models.CacheUserModel]
	sqlDbTodoRepo *sqldb.SqlTodoRepository
	// sqlDbTodoRepo sqldb.ToDoRepositoryInterface
}

type GetUserTodoListResponse struct {
	Id          string
	Description string
	IsComplete  bool
	CreatedAt   time.Time
}

// Interfaces.
type TodoServiceInterface interface {
	GetUserTodoList(userId string) []GetUserTodoListResponse
	AddTodoItem(userId string, description string) (id string)
}

// Functions.
func TodoServiceInit(c *cache.Cache[string, models.CacheUserModel], sqldbTodo *sqldb.SqlTodoRepository) *TodoService {
	return &TodoService{
		cache:         c,
		sqlDbTodoRepo: sqldbTodo,
	}
}

func (s *TodoService) GetUserTodoList(userId string) (r []GetUserTodoListResponse) {

	user, found := s.cache.Get(userId)

	// Check if the user is found in cache and if is still valid.
	if found && user.ExpiresAt.After(time.Now()) {
		fmt.Println("User found and valid. Returning it!")
		for k, v := range user.TodoList.List() {
			r = append(r, GetUserTodoListResponse{Id: k, Description: v.Description, CreatedAt: v.CreatedAt, IsComplete: v.IsComplete})
		}
		return r
	}

	rows, _ := s.sqlDbTodoRepo.FindAllByUser(userId)
	fmt.Println("Fetching info From DB")

	user = s.setCacheNewUser(userId)

	for _, item := range rows {

		r = append(r, GetUserTodoListResponse{Id: item.Id, Description: item.Description, CreatedAt: item.CreatedAt, IsComplete: item.IsComplete})

		s.setCacheUserTodo(&user, item.Id, models.CacheTodoItemModel{
			Description: item.Description,
			CreatedAt:   item.CreatedAt,
			IsComplete:  item.IsComplete,
		})

	}

	return r

}

func (s *TodoService) AddTodoItem(userId string, description string) (id string) {

	// Insert into the database.
	sqlId := s.sqlDbTodoRepo.Insert(userId, description)
	fmt.Println("Inserted record with ID:", sqlId)

	// Update the cache.
	user, found := s.cache.Get(userId)
	if !found {
		fmt.Println("User not found. Creating it!")
		user = s.cache.Set(userId, models.CacheUserModel{
			TodoList:  cache.New[string, models.CacheTodoItemModel](),
			ExpiresAt: time.Now().Add(30 * time.Second),
		})
	}

	todoId := uuid.New().String()
	// user.Ttl = time.Now().Add(30 * time.Second)
	user.TodoList.Set(todoId, models.CacheTodoItemModel{
		Description: description,
		CreatedAt:   time.Now(),
		IsComplete:  false,
	})

	// s.cache.Set(userId, user)

	return todoId

}

func (s *TodoService) setCacheNewUser(userId string) models.CacheUserModel {

	return s.cache.Set(userId, models.CacheUserModel{
		TodoList:  cache.New[string, models.CacheTodoItemModel](),
		ExpiresAt: time.Now().Add(30 * time.Second),
	})

}

func (s *TodoService) setCacheUserTodo(userCache *models.CacheUserModel, todoId string, todo models.CacheTodoItemModel) {

	userCache.TodoList.Set(todoId, todo)
	userCache.ExpiresAt = time.Now().Add(30 * time.Second)

}
