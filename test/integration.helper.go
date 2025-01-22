package test

import (
	"database/sql"
	"net/http/httptest"
	"testing"

	"local.com/todo-list-app/internal/setup"
)

func RunTestServer() (*httptest.Server, func(tb testing.TB)) {

	router, db := setup.ServerSetup("TEST")

	RunUp(db)

	return httptest.NewServer(router), func(tb testing.TB) {
		RunDown(db)
		// log.Println("teardown suite")
	}
}

func RunUp(db *sql.DB) {

	_, err := db.Exec(`DROP TABLE IF EXISTS todo;`)
	if err != nil {
		return
	}

	_, err = db.Exec(`
			CREATE TABLE todo (
				id VARCHAR(36) NOT NULL PRIMARY KEY,
				user_id VARCHAR(36) NOT NULL,
				description TEXT NOT NULL,
				is_completed BOOLEAN NOT NULL DEFAULT FALSE,
				rank INT NOT NULL,
				created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
				deleted_at TIMESTAMP(3)
			);`)
	if err != nil {
		return
	}

}

func RunDown(db *sql.DB) {

	_, err := db.Exec("DROP TABLE IF EXISTS todo")
	if err != nil {
		return
	}

	// fmt.Println("shuting down tests")
	db.Close()
}
