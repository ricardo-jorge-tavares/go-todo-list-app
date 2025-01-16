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
	Rank        int          `json:"rank"`
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

// func (r *SqlTodoRepository) FindAll() ([]Todo, error) {

// 	rows, err := r.db.Query("SELECT * FROM todo")
// 	if err != nil {
// 		log.Fatalf("Unable to execute the query. %v", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var records []Todo
// 	for rows.Next() {
// 		var row Todo
// 		err := rows.Scan(&row.Id, &row.UserId, &row.Description, &row.IsComplete, &row.CreatedAt, &row.DeletedAt)
// 		if err != nil {
// 			log.Fatalf("Unable to scan the row. %v", err)
// 			return nil, err
// 		}
// 		records = append(records, row)
// 	}

// 	return records, nil

// }

func (r *SqlTodoRepository) FindAllByUser(userId string) ([]Todo, error) {

	rows, err := r.db.Query("SELECT id, user_id, description, is_complete, rank, created_at, deleted_at FROM todo WHERE user_id=$1 ORDER BY rank", userId)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
		return nil, err
	}
	defer rows.Close()

	var records []Todo
	for rows.Next() {
		var row Todo
		err := rows.Scan(&row.Id, &row.UserId, &row.Description, &row.IsComplete, &row.Rank, &row.CreatedAt, &row.DeletedAt)
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
	err := r.db.QueryRow("SELECT * FROM todo WHERE id=$1", id).Scan(&row.Id, &row.UserId, &row.Description, &row.IsComplete, &row.Rank, &row.CreatedAt, &row.DeletedAt)

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

	err := r.db.QueryRow(`
		INSERT INTO todo (id, user_id, description, is_complete, rank, created_at)
		SELECT $1, $2, $3, $4, (SELECT COUNT(*) + 1 FROM todo WHERE user_id = $5), $6
		RETURNING id`, id, userId, description, false, userId, time.Now()).Scan(&id)

	if err != nil {
		log.Fatal(err)
	}

	return id
}

func (r *SqlTodoRepository) Update(todoId string, description string) {

	_, err := r.db.Exec("UPDATE todo SET description=$1 WHERE id=$2", description, todoId)
	if err != nil {
		fmt.Println("error updating", err)
		log.Fatal(err)
	}

}

func (r *SqlTodoRepository) UpdateUserRank(userId string, todoId string, toRank int) {

	fmt.Println("updateRank", userId, todoId, toRank)

	_, err := r.db.Exec(`
		WITH q_sorted AS (
			SELECT
				id,
				ROW_NUMBER() OVER (ORDER BY rank) AS row_number
			FROM todo
			WHERE user_id = $1
			AND id <> $2
		)
		UPDATE todo t1
		SET rank =
			(CASE
				WHEN p.id = $2 THEN $3
				WHEN s.row_number >= $3 THEN s.row_number + 1
				ELSE s.row_number
			END)
		FROM todo p
		LEFT JOIN q_sorted s ON p.id = s.id
		WHERE t1.id = p.id`, userId, todoId, toRank)

	if err != nil {
		fmt.Println("error updating", err)
		log.Fatal(err)
	}

}
