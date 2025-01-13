package sqldb

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	Id          string    `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	IsComplete  bool      `json:"isComplete"`
}

// ToDoRepositoryInterface
// type ToDoRepositoryInterface interface {
// 	FindAll() ([]Todo, error)
// 	FindById(id string) (Todo, error)
// 	Insert(description string) string
// }

// SqlTodoRepository implements models.UserRepository
type SqlTodoRepository struct {
	db *sql.DB
}

// NewToDoRepository ..
func NewToDoRepository(db *sql.DB) *SqlTodoRepository {
	return &SqlTodoRepository{
		db: db,
	}
}

func (r *SqlTodoRepository) FindAll() ([]Todo, error) {

	rows, err := r.db.Query("SELECT * FROM todo")
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
		return nil, err
	}
	defer rows.Close()

	var records []Todo
	for rows.Next() {
		var row Todo
		err := rows.Scan(&row.Id, &row.Description, &row.CreatedAt, &row.IsComplete)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
			return nil, err
		}
		records = append(records, row)
	}
	return records, nil

}

// FindById ..
func (r *SqlTodoRepository) FindById(id string) (Todo, error) {

	var row = Todo{}
	err := r.db.QueryRow("SELECT * FROM todo WHERE id=$1", id).Scan(&row.Id, &row.Description, &row.CreatedAt, &row.IsComplete)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return row, nil
	case nil:
		return row, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return row, nil
}

func (r *SqlTodoRepository) Insert(description string) string {

	id := uuid.New().String()

	err := r.db.QueryRow("INSERT INTO todo (id, description, created_at, is_complete) VALUES ($1, $2, $3, $4) RETURNING id", id, description, time.Now(), false).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}
