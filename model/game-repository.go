package model

import (
	"gorm.io/gorm"
)

type GameRepository struct {
	DB *gorm.DB
}

func (g GameRepository) FindALl() []Game {
	var games []Game
	g.DB.Joins("INNER JOIN rounds ON rounds.id = games.round_id AND rounds.dont_show_before < now() AND games.time < rounds.dont_show_before").Order("time, hall asc").Find(&games)
	return games

}

func (g GameRepository) FindByGameDay(id string) (games []Game) {
	g.DB.Joins("INNER JOIN rounds ON rounds.id = games.round_id AND rounds.dont_show_before < now() AND games.time < rounds.dont_show_before").Where("rounds.gameday_id=?", id).Order("time, hall asc").Find(&games)
	return games
}

func (g GameRepository) FindGamesByRoundId(id string) []Game {
	var games []Game
	g.DB.Joins("INNER JOIN rounds ON rounds.id = games.round_id AND rounds.dont_show_before < now() AND games.time < rounds.dont_show_before").Where("round_id=?", id).Order("time, hall asc").Find(&games)
	return games
}

func (g GameRepository) FindAllRounds() []Round {
	var rounds []Round
	g.DB.Where("dont_show_before < now()").Order("id").Find(&rounds)
	return rounds
}

func (g GameRepository) FindGameDaysWithRounds() (gamedays []Gameday) {
	g.DB.Preload("Rounds", func(db *gorm.DB) *gorm.DB {
		return db.Where("rounds.dont_show_before <= now()").Order("rounds.id ASC")
	}).Order("id").Find(&gamedays)

	return gamedays
}

func (g GameRepository) FindGameDays() (gamedays []Gameday) {
	g.DB.Order("day").Find(&gamedays)
	return gamedays
}

func (g GameRepository) FindTableRowsByRound(id string) (tablerows []TableRow) {
	g.DB.Where("round_id=?", id).Order("Place").Find(&tablerows)
	return tablerows
}
