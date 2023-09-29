package entities

type UnprocessedURL struct {
	URL string `json:"url"`
}

type ShortURL struct {
	ID     string `json:"-"`
	URL    string `json:"result"`
	Create bool   `json:"-"`
}
