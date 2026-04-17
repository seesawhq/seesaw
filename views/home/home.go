package home

import (
	"html/template"
	"net/http"
)

func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/home/index.html"))
		tmpl.Execute(w, nil)
	}
}

func New(serverMux *http.ServeMux) {
	serverMux.HandleFunc("/", Index())
}
