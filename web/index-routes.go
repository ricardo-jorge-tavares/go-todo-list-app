package web

import (
	"net/http"

	helpers "local.com/todo-list-app/internal/helpers"
)

func IndexRoute(writer http.ResponseWriter, request *http.Request) {

	t, _ := helpers.ParseView("web/views/index/index.html")
	err := t.Execute(writer, nil)
	helpers.CheckError(err)

}
