package model

import "time"

type Team struct {
	ID        string `sql:"type:text;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
