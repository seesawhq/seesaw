package auth

import (
	"fmt"
	"html/template"
	"net/http"
)

func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl := template.Must(template.ParseFiles("templates/layouts/auth.html", "templates/auth/login.html"))
			tmpl.Execute(w, nil)
		}
		if r.Method == http.MethodPost {
			email := r.PostFormValue("email")
			password := r.PostFormValue("password")
			fmt.Println(email, password)
		}
	}
}

func New(serverMux *http.ServeMux) {
	serverMux.HandleFunc("GET /login/", Login())
	serverMux.HandleFunc("POST /login/", Login())
}
