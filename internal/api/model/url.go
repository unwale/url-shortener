package model

type ShortenURLRequest struct {
	URL   string `json:"url"`
	Alias string `json:"alias,omitempty"`
}

type ShortenURLResponse struct {
	ShortURL string `json:"short_url"`
}

type ResolveURLRequest struct {
	ShortURL string `json:"short_url"`
}
