package model

type ShortenURLRequest struct {
	URL string `json:"url"`
}

type ShortenURLResponse struct {
	ShortURL string `json:"short_url"`
}

type ResolveURLRequest struct {
	ShortURL string `json:"short_url"`
}
