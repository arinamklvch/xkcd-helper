package controller

import (
	"encoding/json"
	"net/http"

	"github.com/arinamklvch/xkcd-helper/internal/dto"
)

// @Summary Get last comic
// @Description Returns the last stored XKCD comic as a "number, title, image URL" object.
// @Tags comics
// @Produce json
// @Success 200 {object} dto.SearchComic "Last comic"
// @Failure 500 {string} string "Failed to send JSON"
// @Router /last-comic [get]
func (h *Handler) getLastComic(w http.ResponseWriter, r *http.Request) {
	lastComic, err := h.service.GetLastComic()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	ouputDto := dto.SearchComic{
		Num:    lastComic.Num,
		Title:  lastComic.Title,
		ImgUrl: lastComic.Img,
	}

	err = json.NewEncoder(w).Encode(ouputDto)
	if err != nil {
		http.Error(w, "failed to send JSON", http.StatusInternalServerError)
	}
}
