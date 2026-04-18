package auth

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
)

func Login(templs fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl := template.Must(template.ParseFS(templs, "templates/layouts/auth.html", "templates/auth/login.html"))
			tmpl.Execute(w, nil)
		}
		if r.Method == http.MethodPost {
			email := r.PostFormValue("email")
			password := r.PostFormValue("password")
			fmt.Println(email, password)
		}
	}
}

func New(templs fs.FS, serverMux *http.ServeMux) {
	serverMux.HandleFunc("GET /login/", Login(templs))
	serverMux.HandleFunc("POST /login/", Login(templs))
}
