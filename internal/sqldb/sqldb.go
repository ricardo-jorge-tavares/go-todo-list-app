package sqldb

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// ConnectDB opens a connection to the database
func ConnectDB() *sql.DB {

	db, err := sql.Open("postgres", "postgres://admin:admin@localhost:5432/todo_list_app?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	// defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db

}
