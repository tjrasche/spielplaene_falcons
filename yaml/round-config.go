package yaml

type GameDay struct {
	Rounds []RoundConfig
	Bucket string
	Name   string
}
type RoundConfig struct {
	Worksheet  string
	GameRanges []struct {
		Start CellDef
		End   CellDef
	}
}

type CellDef struct {
	Col string
	Row int
}
