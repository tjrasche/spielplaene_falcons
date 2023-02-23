package model

import (
	"time"
)

type Game struct {
	ID         string `sql:"type:text;primary_key;"`
	Home       Team
	HomeID     string
	Away       Team
	AwayID     string
	Time       time.Time
	AwayScores int
	HomeScores int
	Hall       string
	Round      Round
	RoundID    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
