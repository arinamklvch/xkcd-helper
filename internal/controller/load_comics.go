package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/arinamklvch/xkcd-helper/internal/dto"
	"github.com/arinamklvch/xkcd-helper/internal/usecase"
)

// @Summary Load comics
// @Description Returns XKCD comics as "number, title" strings for the requested range.
// @Tags comics
// @Produce json
// @Param from query int true "Starting comic number"
// @Param to query int true "Ending comic number"
// @Success 200 {array} dto.LoadComic "Loaded comics"
// @Failure 400 {string} string "Invalid request parameters"
// @Failure 500 {string} string "Failed to send JSON"
// @Router /load-comics [get]
func (h *Handler) loadComics(w http.ResponseWriter, r *http.Request) {
	from, err := strconv.Atoi(r.URL.Query().Get("from"))
	if err != nil {
		http.Error(w, "invalid 'from' parameter", http.StatusBadRequest)
		return
	}

	to, err := strconv.Atoi(r.URL.Query().Get("to"))
	if err != nil {
		http.Error(w, "invalid 'to' parameter", http.StatusBadRequest)
		return
	}

	output, err := h.service.LoadComics(dto.LoadComicsInput{
		From: from,
		To:   to,
	})
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidRange) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		h.logger.Error("failed to load comics", "method", r.Method, "path", r.URL.Path, "from", from, "to", to, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	outputDto := []dto.LoadComic{}
	for _, o := range output {
		outputDto = append(outputDto, dto.LoadComic{
			Num:   o.Num,
			Title: o.Title,
		})
	}
	err = json.NewEncoder(w).Encode(outputDto)
	if err != nil {
		h.logger.Error("failed to send load comics response", "method", r.Method, "path", r.URL.Path, "error", err)
		http.Error(w, "failed to send JSON", http.StatusInternalServerError)
	}
}
