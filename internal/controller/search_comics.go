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
		tmpl, err := template.ParseFiles("templates/search.html")
		if err != nil {
			http.Error(w, "failed to parse template", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "failed to render template", http.StatusInternalServerError)
			return
		}
		return
	}

	output, err := h.service.SearchComics(dto.SearchComicsInput{
		Query: query,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	outputDto := []dto.SearchComic{}
	for _, o := range output {
		outputDto = append(outputDto, dto.SearchComic{
			Num:    o.Num,
			Title:  o.Title,
			ImgUrl: o.Img,
		})
	}

	tmpl, err := template.ParseFiles("templates/comics.html")
	if err != nil {
		http.Error(w, "failed to parse template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, searchComicsPageData{
		Query:  query,
		Comics: outputDto,
	})
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}

}
