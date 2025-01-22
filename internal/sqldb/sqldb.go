package sqldb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

// ConnectDB opens a connection to the database
func ConnectDB(environment string) *sql.DB {

	var env = ""
	if environment == "TEST" {
		env = "TEST_"
	} else {
		env = ""
	}

	var (
		DB_HOST     = strings.TrimSpace(os.Getenv(env + "SQL_HOST"))
		DB_PORT     = strings.TrimSpace(os.Getenv(env + "SQL_PORT"))
		DB_USER     = strings.TrimSpace(os.Getenv(env + "SQL_USER"))
		DB_PASSWORD = strings.TrimSpace(os.Getenv(env + "SQL_PASSWORD"))
		DB_NAME     = strings.TrimSpace(os.Getenv(env + "SQL_DBNAME"))
	)

	if DB_HOST == "" || DB_PORT == "" || DB_USER == "" || DB_PASSWORD == "" || DB_NAME == "" {
		log.Fatal("error: missing required environment variables")
	}

	dbConnection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	db, err := sql.Open("postgres", dbConnection)

	if err != nil {
		log.Fatal(err)
	}

	// defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db

}
