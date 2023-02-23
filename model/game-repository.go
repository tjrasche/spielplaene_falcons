package model

import "gorm.io/gorm"

type GameRepository struct {
	DB *gorm.DB
}

func (g GameRepository) FindALl() []Game {
	var games []Game
	g.DB.Joins("INNER JOIN rounds ON rounds.id = games.round_id AND rounds.dont_show_before < now() AND games.time < rounds.dont_show_before").Order("time").Find(&games)
	return games

}
