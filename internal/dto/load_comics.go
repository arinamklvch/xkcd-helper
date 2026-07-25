package dto

type LoadComicsInput struct {
	From int
	To   int
}

// dto.Comic -- сущность для/из handler
type Comic struct {
	Num   int    `json:"num"`
	Title string `json:"description"`
}
