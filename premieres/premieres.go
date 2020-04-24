package premieres

import (
	"fmt"
	"github.com/ynori7/tvshows/streamer"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ynori7/tvshows/config"
)

const premieresUrl = "https://www.metacritic.com/feature/tv-premiere-dates?page=1"
const oneWeek = 7

type PremieresClient struct {
	httpClient   *http.Client
	conf         config.Config
	premieresUrl string
	now          time.Time
}

func NewPremieresClient(conf config.Config) PremieresClient {
	return PremieresClient{
		httpClient:   &http.Client{},
		conf:         conf,
		premieresUrl: premieresUrl,
		now:          time.Now(),
	}
}

func (pc PremieresClient) GetPotentiallyInterestingPremieres(lastProcessedDate string) (*PremiereList, error) {
	// Request the HTML page.
	res, err := pc.httpClient.Get(pc.premieresUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	premiereSet := make(map[string]*Premiere, 0) //used for deduplication

	// Find the new releases
	done := false
	newestDate := ""
	doc.Find(".listtable tr").Each(func(i int, s *goquery.Selection) {
		if done {
			return //stop after we've scanned the last week
		}

		if s.HasClass("sublistbig") { //This is the date headline. Check if this date was within the last week
			dateRaw := s.Find("th").Text()
			dateParts := strings.Split(dateRaw, " / ")
			if newestDate == "" {
				newestDate = dateParts[1]
			}
			if strings.ToLower(dateParts[1]) == strings.ToLower(lastProcessedDate) {
				done = true
				return
			}
		}

		if !s.HasClass("even") {
			return //this is probably a movie (vod=video-on-demand) or just a text note
		}

		premiere := new(Premiere)

		//Get title
		title := s.Find("td.title a").First()
		if link, ok := title.Attr("href"); !ok || strings.Contains(link, "movie") {
			return //we're not interested in movies
		}
		premiere.Title = strings.TrimSpace(title.Text())
		if premiere.Title == "" || premiere.Title == "Trailer" {
			return //if there was no link then it's probably not an interesting show
		}
		if _, ok := premiereSet[premiere.Title]; ok {
			return //this is apparently a duplicate
		}

		//Check if it's a new series
		newImg := s.Find("td.title img[alt=NEW]")
		if newImg != nil && newImg.Nodes != nil {
			premiere.IsNew = true
		}

		//Get genres
		genresRaw := s.Children().Eq(2).Text()
		genreList := strings.Split(genresRaw, "/")
		if !pc.conf.IsInterestingMainGenre(genreList) {
			return //Not an interesting genre
		}
		premiere.Genres = genreList

		//Get streamer
		networkRaw := s.Children().Eq(3)
		premiere.StreamingOption = pc.getStreamer(networkRaw)

		premiereSet[premiere.Title] = premiere
	})

	//turn the map into a list
	premieres := &PremiereList{
		StartDate: lastProcessedDate,
		EndDate: newestDate,
	}
	for _, r := range premiereSet {
		premieres.Premieres = append(premieres.Premieres, *r)
	}

	return premieres, nil
}

func (pc PremieresClient) getStreamer(s *goquery.Selection) streamer.Streamer {
	netflix := s.Find("img[alt=Netflix]")
	if netflix != nil && netflix.Nodes != nil {
		return streamer.Netflix
	} else {
		switch s.Text() {
		case "Prime Video":
			return streamer.Amazon
		case "Disney+":
			return streamer.Disney
		}
	}
	return streamer.None
}
