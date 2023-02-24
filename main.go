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
	// clean up whole schema (it's only a readmodel to excel anyway)
	tx := db.Exec("DROP TABLE IF EXISTS games CASCADE;DROP TABLE IF EXISTS gamedays CASCADE;DROP TABLE IF EXISTS rounds CASCADE;DROP TABLE IF EXISTS team_rounds CASCADE;DROP TABLE IF EXISTS teams CASCADE;DROP TABLE IF EXISTS table_rows CASCADE;")

	if tx.Error != nil {
		panic(tx.Error)
	}
	err = db.AutoMigrate(model.Team{}, model.TableRow{}, model.Round{}, model.Game{}, model.Gameday{})
	if err != nil {
		panic("Automigration of entities failed!")
	}

	// hack to circumnavigate gorm's unflexible automigration
	tx = db.Exec("ALTER TABLE games ALTER COLUMN time TYPE timestamp without time zone;")
	if tx.Error != nil {
		panic(tx.Error)
	}
	cfgProvider, err := yaml.NewConfigProvider("gamedays/")
	if err != nil {
		panic(err)
	}

	for _, gd := range cfgProvider.GamdeDays {
		ticker := time.NewTicker(10000)
		logger.Printf("Updating table for %s", gd.Name)
		go func(day yaml.GameDay) {
			for range ticker.C {
				excelReader := excel.NewReader(os.Getenv("CF_ACCOUNT_ID"), os.Getenv("CF_ACCESSKEY_ID"), os.Getenv("CF_ACCESS_KEY_SECRET"), db, day)
				err := excelReader.UpdateGames()
				if err != nil {
					logger.Printf("Error retrieving reading excel doc!")
					logger.Printf(err.Error())
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
		wd := time.Now().Weekday()
		gamedays := repo.FindGameDays()
		for _, gd := range gamedays {
			if gd.Day.Weekday() == wd {
				c.Redirect(307, c.Request.RequestURI+"gamedays/"+gd.ID)
			}
		}

		c.Redirect(307, c.Request.RequestURI+"gamedays/"+gamedays[0].ID)

	})
	router.GET("/gamedays/:id", func(c *gin.Context) {
		id := c.Param("id")
		games := repo.FindByGameDay(id)
		gamedays := repo.FindGameDaysWithRounds()
		c.HTML(http.StatusOK, "index.html", gin.H{
			"PageTitle": id,
			"Games":     games,
			"Gamedays":  gamedays,
		})
	})
	router.GET("rounds/:id", func(c *gin.Context) {

		id := c.Param("id")
		games := repo.FindGamesByRoundId(id)
		gamedays := repo.FindGameDaysWithRounds()
		table := repo.FindTableRowsByRound(id)
		c.HTML(http.StatusOK, "round.html", gin.H{
			"PageTitle": id,
			"Games":     games,
			"Gamedays":  gamedays,
			"Table":     table,
			"HasTable":  len(table) > 0,
		})
	})

	_ = router.Run()
}
