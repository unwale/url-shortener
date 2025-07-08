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

type ShortUrlStataRequest struct {
	ShortURL string `json:"short_url"`
}

type ShortUrlStatsResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	ClickCount  int    `json:"click_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
