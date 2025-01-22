package main

import (
	"fmt"
	"net/http"

	"local.com/todo-list-app/internal/helpers"
	"local.com/todo-list-app/internal/setup"
)

func main() {

	router, db := setup.ServerSetup("")
	defer db.Close()

	fmt.Println("Server started at http://localhost:8080")
	serverError := http.ListenAndServe(":8080", router)
	helpers.CheckError(serverError)

}
