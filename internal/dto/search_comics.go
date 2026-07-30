package dto

type SearchComicsInput struct {
	Words []string
}

type SearchComic struct {
	Num    int    `json:"num"`
	Title  string `json:"description"`
	ImgUrl string `json:"image_url"`
}
