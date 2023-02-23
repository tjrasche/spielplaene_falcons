package model

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"testing"
)

func TestFindAll(t *testing.T) {
	dsn := "host=" + os.Getenv("POSTGRES_HOST") + " user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=" + os.Getenv("POSTGRES_DATABASE") + " port=" + os.Getenv("POSTGRES_PORT") + " sslmode=disable TimeZone=Europe/Berlin"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	gr := GameRepository{db: db}
	res := gr.findALl()
	for _, game := range res {
		t.Logf(game.HomeID + " vs . " + game.AwayID + " Halle: " + game.Hall)
	}
}
