package yaml

import (
	"gopkg.in/yaml.v3"
	"os"
)

type ConfigProvider struct {
	GamdeDays []GameDay
}

func NewConfigProvider(GameDaysPath string) (*ConfigProvider, error) {
	// find all gamedays
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fullGameDaysPath := wd + "/" + GameDaysPath
	dir, err := os.ReadDir(fullGameDaysPath)
	var pathsToRead []string
	for _, entry := range dir {
		if !entry.IsDir() {
			pathsToRead = append(pathsToRead, fullGameDaysPath+entry.Name())
		}
	}
	var gamedays []GameDay
	for _, filePath := range pathsToRead {
		var gameDay GameDay
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		decoder := yaml.NewDecoder(file)
		err = decoder.Decode(&gameDay)
		if err != nil {
			return nil, err
		}
		gamedays = append(gamedays, gameDay)

	}
	if err != nil {
		return nil, err
	}
	return &ConfigProvider{GamdeDays: gamedays}, nil
}
