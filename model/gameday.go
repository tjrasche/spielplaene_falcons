package model

import "time"

type Gameday struct {
	ID     string `sql:"type:text;primary_key;"`
	Day    time.Time
	Rounds []Round
}
