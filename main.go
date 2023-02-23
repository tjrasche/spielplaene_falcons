package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"rasche-thalhofer.cloud/init/excel"
	"rasche-thalhofer.cloud/init/model"
	"rasche-thalhofer.cloud/init/yaml"
)

func main() {

	dsn := "host=" + os.Getenv("POSTGRES_HOST") + " user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=" + os.Getenv("POSTGRES_DATABASE") + " port=" + os.Getenv("POSTGRES_PORT") + " sslmode=disable TimeZone=Europe/Berlin"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to database!")
	}
	err = db.AutoMigrate(model.Team{}, model.Round{}, model.Game{})
	if err != nil {
		panic("Automigration of entities failed!")
	}
	cfgProvider, err := yaml.NewConfigProvider("gamedays/")
	if err != nil {
		panic(err)
	}

	for _, gd := range cfgProvider.GamdeDays {
		excelReader := excel.NewReader(os.Getenv("CF_ACCOUNT_ID"), os.Getenv("CF_ACCESSKEY_ID"), os.Getenv("CF_ACCESS_KEY_SECRET"), db, gd)
		err := excelReader.UpdateGames()
		if err != nil {
			panic(err)
		}
	}

}
