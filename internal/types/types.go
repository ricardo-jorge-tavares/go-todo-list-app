package types

import "time"

type TodoListItemType struct {
	Description string
	CreatedAt   time.Time
	IsComplete  bool
}
