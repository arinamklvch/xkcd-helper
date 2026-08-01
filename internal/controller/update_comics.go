package controller

import "net/http"

// @Summary Update comics
// @Description Downloads and stores new XKCD comics that are not yet present in the database.
// @Tags comics
// @Success 200 {string} string "Comics updated"
// @Failure 500 {string} string "Failed to update comics"
// @Router /update [put]
func (h *Handler) updateComics(w http.ResponseWriter, r *http.Request) {
	err := h.service.UpdateComics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}
