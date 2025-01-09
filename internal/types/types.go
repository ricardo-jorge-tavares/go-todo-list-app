package types

import "time"

type TodoListItemType struct {
	Descrition string
	CreatedAt  time.Time
	IsComplete bool
}
