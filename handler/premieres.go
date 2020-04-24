package handler

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ynori7/tvshows/config"
	"github.com/ynori7/tvshows/filter"
	"github.com/ynori7/tvshows/premieres"
	"github.com/ynori7/tvshows/tvshow"
	"github.com/ynori7/tvshows/view"
)

const lastProcessedFile = "lastprocessed.dat"
const defaultDays = time.Duration(7)

type PremieresHandler struct {
	conf            config.Config
	premieresClient premieres.PremieresClient
}

func NewPremieresHandler(
	conf config.Config,
	premieresClient premieres.PremieresClient,
) PremieresHandler {
	return PremieresHandler{
		conf:            conf,
		premieresClient: premieresClient,
	}
}

func (h PremieresHandler) GenerateNewReleasesReport() (*PremieresReport, error) {
	logger := log.WithFields(log.Fields{"Logger": "GenerateNewReleasesReport"})

	lastProcessedDate := h.getLastProcessedDate()

	//Get the premieresList of new premieres
	premieresList, err := h.premieresClient.GetPotentiallyInterestingPremieres(lastProcessedDate)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Error("Error getting new premieres")
		return nil, err
	}

	//Fetch the tv show details and filter
	filterer := filter.NewFilterer(h.conf, tvshow.NewTvShowClient(h.conf), premieresList)
	interestingSeries := filterer.FilterAndEnrich()

	if len(interestingSeries) == 0 {
		return nil, fmt.Errorf("no new series")
	}

	//Split the new and returning series
	newSeries := make([]tvshow.TvShow, 0)
	returningSeries := make([]tvshow.TvShow, 0)
	for _, series := range interestingSeries {
		if series.IsNewSeries {
			newSeries = append(newSeries, series)
		} else {
			returningSeries = append(returningSeries, series)
		}
	}

	//Build HTML output
	template := view.NewHtmlTemplate(newSeries, returningSeries)
	out, err := template.ExecuteHtmlTemplate()
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Error("Error generating html")
		return nil, err
	}

	//Save HTML output to file
	dateString := time.Now().Format("20060102") //yyyyMMdd
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s-%s.html", config.CliConf.OutputPath, h.conf.Title, dateString), []byte(out), 0644)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Warn("Error saving html to file")
		return nil, err
	}

	//Mark where we left off
	if err := h.updateLastProcessedDate(premieresList.EndDate); err != nil {
		logger.WithFields(log.Fields{"error": err}).Warn("Error updating last processed date")
	}

	return &PremieresReport{
		Html:      out,
		StartDate: premieresList.StartDate,
		EndDate:   premieresList.EndDate,
	}, nil
}

func (h PremieresHandler) getLastProcessedDate() string {
	dat, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", config.CliConf.LastProcessedPath, lastProcessedFile))
	if err != nil || len(strings.TrimSpace(string(dat))) == 0 {
		lastWeek := time.Now().Add(-1 * defaultDays * 24 * time.Hour)
		return fmt.Sprintf("%s %d", lastWeek.Month(), lastWeek.Day())
	}

	return strings.TrimSpace(string(dat))
}

func (h PremieresHandler) updateLastProcessedDate(date string) error {
	return ioutil.WriteFile(fmt.Sprintf("%s/%s", config.CliConf.LastProcessedPath, lastProcessedFile), []byte(date), 0644)
}
