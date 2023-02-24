package yaml

import "time"

type GameDay struct {
	Rounds []RoundConfig
	Bucket string
	Name   string
	Day    time.Time
}
type RoundConfig struct {
	DontShowBefore time.Time
	Worksheet      string
	GameRanges     []struct {
		Start CellDef
		End   CellDef
	}
	TableRanges []struct {
		Start CellDef
		End   CellDef
	}
}

type CellDef struct {
	Col string
	Row int
}
