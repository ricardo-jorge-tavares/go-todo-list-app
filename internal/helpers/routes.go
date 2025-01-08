package helpers

import (
	"text/template"

	config "local.com/todo-list-app/internal/config"
)

func ParseView(view string) (*template.Template, error) {

	return template.ParseFiles(append(config.UiLayoutViews(), view)...)
	// if err != nil {
	// 	return nil, err
	// }
	// return t, nil
}
