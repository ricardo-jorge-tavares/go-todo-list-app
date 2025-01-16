package models

import (
	"time"

	"local.com/todo-list-app/internal/cache"
)

type CacheUserModel struct {
	TodoList  *cache.Cache[string, CacheTodoItemModel]
	ExpiresAt time.Time
}

type CacheTodoItemModel struct {
	Description string
	IsComplete  bool
	Rank        int
	CreatedAt   time.Time
}
