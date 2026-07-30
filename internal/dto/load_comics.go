package dto

type LoadComicsInput struct {
	From int
	To   int
}

// dto.LoadComic -- сущность для/из handler
type LoadComic struct {
	Num   int    `json:"num"`
	Title string `json:"description"`
}
