package model

import "gorm.io/gorm"

type GameRepository struct {
	DB *gorm.DB
}

func (g GameRepository) FindALl() []Game {
	var games []Game
	g.DB.Joins("INNER JOIN rounds ON rounds.id = games.round_id AND rounds.dont_show_before < now() AND games.time < rounds.dont_show_before").Order("time, hall asc").Find(&games)
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
