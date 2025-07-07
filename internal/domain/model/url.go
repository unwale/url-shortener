package model

type URL struct {
	ID        string
	Original  string
	Shortened string
	CreatedAt string
	UpdatedAt string
}

type Error struct {
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}
