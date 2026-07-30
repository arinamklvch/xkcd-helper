package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/arinamklvch/xkcd-helper/internal/dto"
)

func (h *Handler) searchComics(w http.ResponseWriter, r *http.Request) {
	words := strings.Fields(r.URL.Query().Get("q"))
	output, err := h.service.SearchComics(dto.SearchComicsInput{
		Words: words,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	ouputDto := []dto.SearchComic{}
	for _, o := range output {
		ouputDto = append(ouputDto, dto.SearchComic{
			Num:    o.Num,
			Title:  o.Title,
			ImgUrl: o.Img,
		})
	}
	err = json.NewEncoder(w).Encode(ouputDto)
	if err != nil {
		http.Error(w, "failed to send JSON", http.StatusInternalServerError)
	}
}
