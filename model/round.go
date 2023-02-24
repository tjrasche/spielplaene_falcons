package model

import "time"

type Round struct {
	ID             string `sql:"type:text;primary_key;"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DontShowBefore time.Time
	Games          []Game
	Gameday        Gameday
	GamedayID      string
}
