package postgresql

import (
	"database/sql"
	"time"
)

type Product struct {
	ID          int32          `json:"id"`
	Name        string         `json:"name"`
	Category    string         `json:"category"`
	Price       int32          `json:"price"`
	IsAvailable sql.NullBool   `json:"is_available"`
	CreatedAt   sql.NullTime   `json:"created_at"`
}
