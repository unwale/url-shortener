// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Url struct {
	ID          int32
	OriginalUrl string
	ShortUrl    string
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
	ClickCount  int64
}
