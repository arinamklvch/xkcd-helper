package controller

import (
	"html/template"
	"net/http"
)

func (h *Handler) webLogin(w http.ResponseWriter, r *http.Request) {
	// [sign in] -- отправляется POST /login
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "failed to parse template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}
