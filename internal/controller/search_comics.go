package controller

import (
	"html/template"
	"net/http"

	"github.com/arinamklvch/xkcd-helper/internal/dto"
)

type searchComicsPageData struct {
	Query  string
	Comics []dto.SearchComic
}

// @Summary Search comics
// @Description Returns XKCD comics as "number, title, image URL" objects for the requested search query.
// @Tags comics
// @Produce json
// @Param q query string true "Search query words separated by spaces"
// @Success 200 {array} dto.SearchComic "Matching comics"
// @Failure 500 {string} string "Failed to send JSON"
// @Router /search-comics [get]
func (h *Handler) searchComics(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	// no query parameters
	if query == "" {
		tmpl := template.Must(template.ParseFiles("templates/search.html"))
		tmpl.Execute(w, nil)
		return
	}

	output, err := h.service.SearchComics(dto.SearchComicsInput{
		Query: query,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ouputDto := []dto.SearchComic{}
	for _, o := range output {
		ouputDto = append(ouputDto, dto.SearchComic{
			Num:    o.Num,
			Title:  o.Title,
			ImgUrl: o.Img,
		})
	}

	tmpl := template.Must(template.ParseFiles("templates/comics.html"))
	tmpl.Execute(w, searchComicsPageData{
		Query:  query,
		Comics: ouputDto,
	})

}
