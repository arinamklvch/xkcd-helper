package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/arinamklvch/xkcd-helper/internal/dto"
	"github.com/arinamklvch/xkcd-helper/internal/usecase"
)

type Handler struct {
	service *usecase.Service
}

func NewHandler(service *usecase.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// handler loads XKCD comics in the requested numeric range.
//
// @Summary Load comics
// @Description Returns XKCD comics as "number, title" strings for the requested range.
// @Tags comics
// @Produce json
// @Param from query int true "Starting comic number"
// @Param to query int true "Ending comic number"
// @Success 200 {array} dto.Comic "Loaded comics"
// @Failure 500 {string} string "Failed to send JSON"
// @Router /load-comics [get]
func (h *Handler) LoadComics(w http.ResponseWriter, r *http.Request) {
	from, err := strconv.Atoi(r.URL.Query().Get("from"))
	if err != nil {
		http.Error(w, "Invalid from", http.StatusBadRequest)
		return
	}

	to, err := strconv.Atoi(r.URL.Query().Get("to"))
	if err != nil {
		http.Error(w, "Invalid to", http.StatusBadRequest)
		return
	}

	output, err := h.service.LoadComics(dto.LoadComicsInput{
		From: from,
		To:   to,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	ouputDTO := []dto.Comic{}
	for _, o := range output {
		ouputDTO = append(ouputDTO, dto.Comic{
			Num:   o.Num,
			Title: o.Title,
		})
	}
	err = json.NewEncoder(w).Encode(ouputDTO)
	if err != nil {
		http.Error(w, "Failed to send JSON", http.StatusInternalServerError)
	}
}
