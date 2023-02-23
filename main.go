package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"rasche-thalhofer.cloud/init/excel"
	"rasche-thalhofer.cloud/init/model"
	"rasche-thalhofer.cloud/init/yaml"
	"time"
)

func main() {
	logger := log.Default()
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to database!")
	}
	err = db.AutoMigrate(model.Team{}, model.Round{}, model.Game{})
	if err != nil {
		panic("Automigration of entities failed!")
	}

	// hack to circumnavigate gorm's unflexible automigration
	db.Exec("ALTER TABLE games ALTER COLUMN time TYPE timestamp without time zone;")
	cfgProvider, err := yaml.NewConfigProvider("gamedays/")
	if err != nil {
		panic(err)
	}

	for _, gd := range cfgProvider.GamdeDays {
		ticker := time.NewTicker(10000)
		go func(day yaml.GameDay) {
			for range ticker.C {
				excelReader := excel.NewReader(os.Getenv("CF_ACCOUNT_ID"), os.Getenv("CF_ACCESSKEY_ID"), os.Getenv("CF_ACCESS_KEY_SECRET"), db, day)
				err := excelReader.UpdateGames()
				if err != nil {
					logger.Printf("Error retrieving reading excel doc!")
				}
			}
		}(gd)

	}

	// setup webserver
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	repo := model.GameRepository{
		DB: db,
	}
	router.GET("/", func(c *gin.Context) {
		games := repo.FindALl()
		rounds := repo.FindAllRounds()
		c.HTML(http.StatusOK, "index.html", gin.H{
			"PageTitle": "Alle Spiele",
			"Games":     games,
			"Rounds":    rounds,
		})
	})
	router.GET("rounds/:id", func(c *gin.Context) {

		id := c.Param("id")
		games := repo.FindGamesByRoundId(id)
		rounds := repo.FindAllRounds()

		c.HTML(http.StatusOK, "index.html", gin.H{
			"PageTitle": "Spiele " + id,
			"Games":     games,
			"Rounds":    rounds,
		})
	})

	_ = router.Run()
}
