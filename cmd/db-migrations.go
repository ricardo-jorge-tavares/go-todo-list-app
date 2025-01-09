package main

import (
	"database/sql"

	"local.com/todo-list-app/internal/helpers"
	"local.com/todo-list-app/internal/sqldb"
)

func createTodoTable(db *sql.DB) {

	_, err := db.Exec(`DROP TABLE IF EXISTS todo;`)
	helpers.CheckError(err)

	_, err = db.Exec(`
		CREATE TABLE todo (
			id VARCHAR(36) NOT NULL PRIMARY KEY,
			description TEXT NOT NULL,
			created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
			is_complete BOOLEAN NOT NULL DEFAULT FALSE
		);`)
	helpers.CheckError(err)

}

func main() {

	db := sqldb.ConnectDB()

	createTodoTable(db)

}
