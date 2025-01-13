package services

import (
	"fmt"
	"time"

	"local.com/todo-list-app/internal/cache"
	"local.com/todo-list-app/internal/sqldb"
	"local.com/todo-list-app/internal/types"
)

// Structs.
type TodoService struct {
	cache         *cache.Cache[string, types.TodoListItemType]
	sqlDbTodoRepo *sqldb.SqlTodoRepository
	// sqlDbTodoRepo sqldb.ToDoRepositoryInterface
}

type GetTodoListResponse struct {
	Id          string
	Description string
}

// Interfaces.
type TodoServiceInterface interface {
	GetTodoList(userId string) []GetTodoListResponse
}

// Functions.
func TodoServiceInit(c *cache.Cache[string, types.TodoListItemType], sqldbTodo *sqldb.SqlTodoRepository) *TodoService {
	return &TodoService{
		cache:         c,
		sqlDbTodoRepo: sqldbTodo,
	}
}

func (s *TodoService) GetTodoList(userId string) (r []GetTodoListResponse) {

	s.cache.Set(userId, types.TodoListItemType{Description: "Go to the gym", CreatedAt: time.Now(), IsComplete: false})

	for k, v := range s.cache.List() {
		fmt.Printf("Key: %s, Value: %s | %s | %v\n", k, v.Description, v.CreatedAt, v.IsComplete)
	}

	sqlId := s.sqlDbTodoRepo.Insert("Go to the gym " + userId)
	fmt.Println("Inserted record with ID:", sqlId)

	row, _ := s.sqlDbTodoRepo.FindAll()
	for _, todo := range row {
		r = append(r, GetTodoListResponse{todo.Id, todo.Description})
		fmt.Println("From DB", todo.Id, todo.Description)
	}

	return r
}
