package models

import (
	"database/sql"
)

type Todo struct {
	Id			int
	Title		string
	Description string
	CompletedAt	sql.NullString
	DeletedAt	sql.NullString
}