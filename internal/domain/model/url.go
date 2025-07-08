package model

type Url struct {
	OriginalUrl string
	ShortUrl    string
	ClickCount  int64
	CreatedAt   string
	UpdatedAt   string
}

type Error struct {
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}
