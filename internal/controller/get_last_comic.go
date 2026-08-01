package controller

import (
	"encoding/json"
	"net/http"

	"github.com/arinamklvch/xkcd-helper/internal/dto"
)

// @Summary Get last comic
// @Description Returns the latest stored XKCD comic as a "number, title, image URL" object.
// @Tags comics
// @Produce json
// @Success 200 {object} dto.SearchComic "Latest comic"
// @Failure 500 {string} string "Failed to send JSON"
// @Router /last-comics [get]
func (h *Handler) getLastComic(w http.ResponseWriter, r *http.Request) {
	latestComic, err := h.service.GetLatestComic()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	ouputDto := dto.SearchComic{
		Num:    latestComic.Num,
		Title:  latestComic.Title,
		ImgUrl: latestComic.Img,
	}

	err = json.NewEncoder(w).Encode(ouputDto)
	if err != nil {
		http.Error(w, "failed to send JSON", http.StatusInternalServerError)
	}
}
