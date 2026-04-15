package home

import "net/http"

func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>Hello world</h1>"))
	}
}

func New(serverMux *http.ServeMux) {
	serverMux.HandleFunc("/", Index())
}
