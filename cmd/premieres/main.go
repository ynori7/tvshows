package main

import (
	"fmt"
	"github.com/ynori7/tvshows/handler"
	"github.com/ynori7/tvshows/tvshow"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/ynori7/tvshows/config"
	"github.com/ynori7/tvshows/premieres"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)
	logger := log.WithFields(log.Fields{"Logger": "main"})

	//Get the cli flags
	config.ParseCliFlags()
	if config.CliConf.ConfigFile == "" {
		logger.Fatal("You must specify the path to the config file")
	}

	//Get the config
	data, err := ioutil.ReadFile(config.CliConf.ConfigFile)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error reading config file")
	}

	var conf config.Config
	if err := conf.Parse(data); err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error parsing config")
	}

	premieresHandler := handler.NewPremieresHandler(conf, premieres.NewPremieresClient(conf))
	newPremieres, err := premieresHandler.GetNewPremieres()
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Error("Error getting interesting new premieres")
		return
	}

	fmt.Println(newPremieres.StartDate, " - ", newPremieres.EndDate)

	tvShowClient := tvshow.NewTvShowClient(conf)
	for _, premiere := range newPremieres.Premieres {
		fmt.Println(premiere)
		fmt.Println(tvShowClient.SearchForTvSeriesTitle(premiere.Title))
	}
}
