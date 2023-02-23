package main

import (
	"fmt"
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
	router.LoadHTMLGlob(os.Getenv("KO_DATA_PATH") + "/*")
	repo := model.GameRepository{
		DB: db,
	}
	router.GET("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "index.html", gin.H{
			"PageTitle": "My Page",
			"Games":     repo.FindALl(),
		})
	})
	router.GET("rounds/:id", func(c *gin.Context) {

		id := c.Param("id")

		c.HTML(http.StatusOK, "index.html", gin.H{
			"PageTitle": "My Page",
			"Games":     id,
		})
	})

	_ = router.Run()
}

func getData(id string) string {
	// Your logic to fetch data based on the ID goes here
	return fmt.Sprintf("Data for ID %s", id)
}
