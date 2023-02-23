package model

import "time"

type Team struct {
	ID        string  `sql:"type:text;primary_key;"`
	Rounds    []Round `gorm:"many2many:team_rounds;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
