package controller

import (
	"html/template"
	"net/http"
)

func (h *Handler) webLogin(w http.ResponseWriter, r *http.Request) {
	// [sign in] -- отправляется POST /login
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, nil)
}
