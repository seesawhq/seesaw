package controllers

import (
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
)

func Index(logger *slog.Logger, templs fs.FS) http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(templs, "templates/layouts/auth.html", "templates/auth/login.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}

func Complete(logger *slog.Logger, templs fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFS(templs, "templates/layouts/auth.html", "templates/auth/login.html"))
		tmpl.Execute(w, nil)
	}
}

func New(logger *slog.Logger, templs fs.FS, serverMux *http.ServeMux) {
	serverMux.HandleFunc("/setup", Index(logger, templs))
	serverMux.HandleFunc("/setup/complete", Complete(logger, templs))
}
