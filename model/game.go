package model

import (
	"gorm.io/gorm"
	"time"
)

type Game struct {
	gorm.Model
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
}
