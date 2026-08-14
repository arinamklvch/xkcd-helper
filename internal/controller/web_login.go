package controller

import (
	"html/template"
	"net/http"
)

func (h *Handler) webLogin(w http.ResponseWriter, r *http.Request) {
	// [sign in] -- отправляется POST /login
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		h.logger.Error("failed to parse login template", "method", r.Method, "path", r.URL.Path, "error", err)
		http.Error(w, "failed to parse template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		h.logger.Error("failed to render login template", "method", r.Method, "path", r.URL.Path, "error", err)
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}
