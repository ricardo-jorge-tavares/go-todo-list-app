package sqldb

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	Id          string       `json:"id"`
	UserId      string       `json:"userId"`
	Description string       `json:"description"`
	IsComplete  bool         `json:"isComplete"`
	CreatedAt   time.Time    `json:"createdAt"`
	DeletedAt   sql.NullTime `json:"deletedAt"`
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
		err := rows.Scan(&row.Id, &row.UserId, &row.Description, &row.IsComplete, &row.CreatedAt, &row.DeletedAt)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
			return nil, err
		}
		records = append(records, row)
	}

	return records, nil

}

func (r *SqlTodoRepository) FindAllByUser(userId string) ([]Todo, error) {

	rows, err := r.db.Query("SELECT * FROM todo WHERE user_id=$1", userId)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
		return nil, err
	}
	defer rows.Close()

	var records []Todo
	for rows.Next() {
		var row Todo
		err := rows.Scan(&row.Id, &row.UserId, &row.Description, &row.IsComplete, &row.CreatedAt, &row.DeletedAt)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
			return nil, err
		}
		records = append(records, row)
	}
	return records, nil

}

// FindById ..
func (r *SqlTodoRepository) FindByTodoId(id string) (Todo, error) {

	var row = Todo{}
	err := r.db.QueryRow("SELECT * FROM todo WHERE id=$1", id).Scan(&row.Id, &row.UserId, &row.Description, &row.IsComplete, &row.CreatedAt, &row.DeletedAt)

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

func (r *SqlTodoRepository) Insert(userId string, description string) string {

	id := uuid.New().String()

	err := r.db.QueryRow("INSERT INTO todo (id, user_id, description, is_complete, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id", id, userId, description, false, time.Now()).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}

	return id
}
