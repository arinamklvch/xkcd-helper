package controller

import (
	"net/http"
)

// @Summary Update comics
// @Description Downloads and stores new XKCD comics that are not yet present in the database.
// @Tags comics
// @Success 200 {string} string "Comics updated"
// @Failure 500 {string} string "Failed to update comics"
// @Router /update [put]
func (h *Handler) updateComics(w http.ResponseWriter, r *http.Request) {
	err := h.service.UpdateComics()
	if err != nil {
		h.logger.Error("failed to update comics", "method", r.Method, "path", r.URL.Path, "error", err)
		http.Error(w, "failed to update comics", http.StatusInternalServerError)
		return
	}
}
