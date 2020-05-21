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

		//Check if it's a movie
		titleLink := s.Find("td.title a").First()
		if link, ok := titleLink.Attr("href"); ok && strings.Contains(link, "movie") {
			return //we're not interested in movies
		}
		movieFlag := s.Find("td.title img[alt=MOVIE]")
		if movieFlag != nil && movieFlag.Nodes != nil {
			return //this is a movie
		}

		//Get the title text
		titleRaw := s.Find("td.title").Text()
		premiere.Title = pc.cleanTitle(titleRaw)
		if premiere.Title == ""  {
			return //if there was no title then this is a garbage entry
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

func (pc PremieresClient) cleanTitle(t string) string {
	t = strings.ReplaceAll(t, "Trailer2", "")
	t = strings.ReplaceAll(t, "Trailer", "")
	t = strings.ReplaceAll(t, "Opening scene", "")
	t = strings.ReplaceAll(t, "Full 1st episode", "")
	t = strings.ReplaceAll(t, "Red-band trailer", "")
	return strings.TrimSpace(t)
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
