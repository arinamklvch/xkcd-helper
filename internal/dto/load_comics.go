package dto

type LoadComicsInput struct {
	From int
	To   int
}

type Comic struct {
	Num   int    `json:"num"`
	Title string `json:"description"`
}
