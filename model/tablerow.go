package model

type TableRow struct {
	ID      string
	Team    Team
	TeamID  string
	Round   Round
	RoundId string
	Wins,
	Diff,
	Place int
}
