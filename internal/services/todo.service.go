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
	sqlDbTodoRepo *sqldb.SqlTodoRepository
	// sqlDbTodoRepo sqldb.ToDoRepositoryInterface
}

type GetUserTodoListResponse struct {
	Id          string
	Description string
	IsComplete  bool
	Rank        int
	CreatedAt   time.Time
}

// Interfaces.
type TodoServiceInterface interface {
	GetUserTodoList(userId string) []GetUserTodoListResponse
	AddTodoItem(userId string, description string) (id string)
	UpdateTodoItemDescription(userId string, todoId string, description string) (id string)
	UpdateTodoItemRank(userId string, todoId string, rank int) (id string)
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
			r = append(r, GetUserTodoListResponse{Id: k, Description: v.Description, IsComplete: v.IsComplete, Rank: v.Rank, CreatedAt: v.CreatedAt})
		}

		sort.Slice(r, func(i, j int) bool {
			return r[i].Rank < r[j].Rank
		})

		return r
	}

	fmt.Println("Fetching info From DB")
	rows, _ := s.sqlDbTodoRepo.FindAllByUser(userId)

	user = s.cacheSetUser(userId)

	for _, item := range rows {

		user.TodoList.Set(item.Id, models.CacheTodoItemModel{
			Description: item.Description,
			IsComplete:  item.IsComplete,
			Rank:        item.Rank,
			CreatedAt:   item.CreatedAt,
		})

		r = append(r, GetUserTodoListResponse{Id: item.Id, Description: item.Description, IsComplete: item.IsComplete, Rank: item.Rank, CreatedAt: item.CreatedAt})

	}

	return r

}

func (s *TodoService) AddTodoItem(userId string, description string) (id string) {

	// Insert into the database.
	sqlTodoId := s.sqlDbTodoRepo.Insert(userId, description)
	fmt.Println("Inserted record with ID:", sqlTodoId)

	s.cacheInvalidateUser(userId)

	return sqlTodoId

}

func (s *TodoService) UpdateTodoItemDescription(userId string, todoId string, description string) (id string) {

	// Update the database.
	s.sqlDbTodoRepo.Update(todoId, description)

	s.cacheInvalidateUser(userId)

	return todoId

}

func (s *TodoService) UpdateTodoItemRank(userId string, todoId string, rank int) (id string) {

	// Update the database.
	s.sqlDbTodoRepo.UpdateUserRank(userId, todoId, rank)

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
