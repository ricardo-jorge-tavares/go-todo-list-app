package main

import (
	"database/sql"

	"github.com/joho/godotenv"
	"local.com/todo-list-app/internal/helpers"
	"local.com/todo-list-app/internal/sqldb"
)

func createTodoTable(db *sql.DB) {

	_, err := db.Exec(`DROP TABLE IF EXISTS todo;`)
	helpers.CheckError(err)

	_, err = db.Exec(`
		CREATE TABLE todo (
			id VARCHAR(36) NOT NULL PRIMARY KEY,
			user_id VARCHAR(36) NULL,
			description TEXT NOT NULL,
			is_complete BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP(3)
		);`)
	helpers.CheckError(err)

}

func main() {

	err := godotenv.Load()
	helpers.CheckError(err)

	db := sqldb.ConnectDB()

	createTodoTable(db)

}
