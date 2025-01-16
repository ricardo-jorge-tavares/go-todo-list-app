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
	IsCompleted bool
	Rank        int
	CreatedAt   time.Time
}
