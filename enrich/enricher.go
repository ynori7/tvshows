package enrich

import (
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/ynori7/tvshows/config"
	"github.com/ynori7/tvshows/premieres"
	"github.com/ynori7/tvshows/tvshow"
	"github.com/ynori7/workerpool"
)

type Enricher struct {
	conf               config.Config
	potentialPremieres *premieres.PremiereList
	tvshowClient       tvshow.ImdbClient
}

func NewEnricher(conf config.Config, discographyClient tvshow.ImdbClient, premieres *premieres.PremiereList) Enricher {
	return Enricher{
		conf:               conf,
		potentialPremieres: premieres,
		tvshowClient:       discographyClient,
	}
}

func (f Enricher) FilterAndEnrich() []tvshow.TvShow {
	logger := log.WithFields(log.Fields{"Logger": "FilterAndEnrich"})

	//Process results
	series := make([]tvshow.TvShow, 0)

	//Set up worker pool
	workerPool := workerpool.NewWorkerPool(5,
		func(result interface{}) {
			r := result.(*tvshow.TvShow)
			logger.WithFields(log.Fields{"Artist": r.Title}).Debug("Found interesting series")
			series = append(series, *r)
		},
		func(err error) {
			logger.WithFields(log.Fields{"error": err}).Error("Error looking up series data")
		},
		f.processPremiere,
	)

	//Do the work
	if err := workerPool.Work(f.potentialPremieres.Premieres); err != nil {
		logger.WithFields(log.Fields{"error": err}).Error("Error processing jobs")
	}

	//Sort the results
	sort.Slice(series, func(i, j int) bool {
		return series[i].Score > series[j].Score
	})

	return series
}

func (f Enricher) processPremiere(job interface{}) (result interface{}, err error) {
	j := job.(premieres.Premiere)

	imdbLink, err := f.tvshowClient.SearchForTvSeriesTitle(j.Title)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, j.Title)
	}

	series, err := f.tvshowClient.GetTvShowData(imdbLink)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, j.Title)
	}
	series.IsNewSeries = j.IsNew
	series.StreamingOption = j.StreamingOption

	return series, nil
}


