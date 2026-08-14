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
// @Description Renders a search form or XKCD comics matching the requested search query.
// @Tags comics
// @Produce html
// @Param q query string false "Search query words separated by spaces"
// @Success 200 {string} string "Search page or matching comics page"
// @Failure 500 {string} string "Internal server error"
// @Router /search-comics [get]
func (h *Handler) searchComics(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	// no query parameters
	if query == "" {
		tmpl, err := template.ParseFiles("templates/search.html")
		if err != nil {
			h.logger.Error("failed to parse search template", "method", r.Method, "path", r.URL.Path, "error", err)
			http.Error(w, "failed to parse template", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			h.logger.Error("failed to render search template", "method", r.Method, "path", r.URL.Path, "error", err)
			http.Error(w, "failed to render template", http.StatusInternalServerError)
			return
		}
		return
	}

	output, err := h.service.SearchComics(dto.SearchComicsInput{
		Query: query,
	})
	if err != nil {
		h.logger.Error("failed to search comics", "method", r.Method, "path", r.URL.Path, "query", query, "error", err)
		http.Error(w, "failed to search comics", http.StatusInternalServerError)
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
		h.logger.Error("failed to parse comics template", "method", r.Method, "path", r.URL.Path, "error", err)
		http.Error(w, "failed to parse template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, searchComicsPageData{
		Query:  query,
		Comics: outputDto,
	})
	if err != nil {
		h.logger.Error("failed to render comics template", "method", r.Method, "path", r.URL.Path, "query", query, "error", err)
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}
