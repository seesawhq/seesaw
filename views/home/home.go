package home

import (
	"html/template"
	"io/fs"
	"net/http"
)

func Index(templs fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFS(templs, "templates/home/index.html"))
		tmpl.Execute(w, nil)
	}
}

func New(templs fs.FS, serverMux *http.ServeMux) {
	serverMux.HandleFunc("/", Index(templs))
}
