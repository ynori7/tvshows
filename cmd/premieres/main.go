package main

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/ynori7/tvshows/application"
	"github.com/ynori7/tvshows/config"
	"github.com/ynori7/tvshows/email"
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

	premieresReporter := application.NewPremieresReporter(conf, premieres.NewPremieresClient(conf))
	newPremieresReport, err := premieresReporter.GeneratePremieresReport()
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Error("Error getting interesting new premieres")
		return
	}

	if conf.Email.Enabled {
		mailer := email.NewMailer(conf)
		if err := mailer.SendMail(email.GetNewReleasesSubjectLine(newPremieresReport.StartDate, newPremieresReport.EndDate), newPremieresReport.Html); err != nil {
			logger.WithFields(log.Fields{"error": err}).Error("Error sending email")
		}
	}
}
