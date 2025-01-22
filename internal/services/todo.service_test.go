package services_test

import (
	"database/sql"
	"testing"
	"time"

	"local.com/todo-list-app/internal/cache"
	"local.com/todo-list-app/internal/models"
	"local.com/todo-list-app/internal/services"
	"local.com/todo-list-app/internal/sqldb"
)

type MockTodoRepository struct {
	FindUserItemsFunc       func(userId string) ([]sqldb.Todo, error)
	FindUserItemsFuncCalled int
	InsertItemFunc          func(userId string, description string) string
}

func (s *MockTodoRepository) FindUserItems(userId string) ([]sqldb.Todo, error) {
	s.FindUserItemsFuncCalled++
	return s.FindUserItemsFunc(userId)
}
func (s *MockTodoRepository) InsertItem(userId string, description string) string {
	return s.InsertItemFunc(userId, description)
}
func (s *MockTodoRepository) UpdateItemDescription(todoId string, description string) {}
func (s *MockTodoRepository) UpdateItemRank(userId string, todoId string, toRank int) {}
func (s *MockTodoRepository) UpdateItemIsCompleted(todoId string)                     {}
func (s *MockTodoRepository) DeleteItem(todoId string)                                {}

func TestServicesSuite(t *testing.T) {

	userCache := cache.New[string, models.CacheUserModel]()
	repo := &MockTodoRepository{}

	userId := "f6645494-56f5-40e4-a3f9-b6ae8935006b"

	t.Run("should return a user from DB and insert on cache", func(t *testing.T) {

		repo = &MockTodoRepository{
			FindUserItemsFunc: func(userId string) ([]sqldb.Todo, error) {
				return []sqldb.Todo{
					{
						Id:          "some-random-id",
						UserId:      userId,
						Description: "A task to do something",
						IsCompleted: false,
						Rank:        1,
						CreatedAt:   time.Now(),
						DeletedAt:   sql.NullTime{Valid: false},
					},
				}, nil
			},
		}

		service := services.TodoServiceInit(userCache, repo)
		returnData := service.GetUserTodoList(userId)

		// FindUserItemsFunc should have been called once.
		if repo.FindUserItemsFuncCalled != 1 {
			t.Errorf("Expected FindUserItemsFunc to be called, got %v", repo.FindUserItemsFuncCalled)
		}

		if (len(returnData) != 1) || (returnData[0].Description != "A task to do something") {
			t.Errorf("Expected returned information to be correct, got %v", returnData)
		}

		_, found := userCache.Get(userId)
		if found == false {
			t.Errorf("Expected user to be on cache")
		}

	})

	t.Run("should return a user from cache and NOT from DB", func(t *testing.T) {

		service := services.TodoServiceInit(userCache, repo)
		returnData := service.GetUserTodoList(userId)

		// FindUserItemsFunc should have been called once (on the previous test).
		if repo.FindUserItemsFuncCalled != 1 {
			t.Errorf("Expected FindUserItemsFunc to be called, got %v", repo.FindUserItemsFuncCalled)
		}

		if (len(returnData) != 1) || (returnData[0].Description != "A task to do something") {
			t.Errorf("Expected returned information to be correct, got %v", returnData)
		}

		_, found := userCache.Get(userId)
		if found == false {
			t.Errorf("Expected user to be on cache")
		}

	})

}
