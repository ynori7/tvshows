package premieres

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ynori7/tvshows/streamer"

	"github.com/PuerkitoBio/goquery"
	"github.com/ynori7/tvshows/config"
)

const premieresUrl = "https://www.metacritic.com/news/tv-calendar-archive-of-past-dates/"
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

	//Look for each new date
	doc.Find(".c-CmsContent h3").Each(func(i int, s *goquery.Selection) {
		if done {
			return //stop after we've scanned the last week
		}

		dateRaw := s.Text()
		dateParts := strings.Split(dateRaw, " / ")
		if newestDate == "" {
			newestDate = strings.TrimSpace(dateParts[1])
		}
		if strings.ToLower(strings.TrimSpace(dateParts[1])) == strings.ToLower(lastProcessedDate) {
			done = true
			return
		}

		//Find the list of premieres for this date
		s.Next().Find("tr").Each(func(i int, s *goquery.Selection) {
			premiere := new(Premiere)

			//Check if it's a movie
			titleLink := s.Find("td:nth-child(2) a").First()
			if link, ok := titleLink.Attr("href"); ok && strings.Contains(link, "movie") {
				return //we're not interested in movies
			}
			movieFlag := s.Find("td:nth-child(2) img[alt=movie]")
			if movieFlag != nil && movieFlag.Nodes != nil {
				return //this is a movie
			}

			//Get the title text
			titleRaw := s.Find("td:nth-child(2) strong").First().Text()
			if titleRaw == "" || titleRaw == "($)" { //sometimes it's a link and sometimes it's a strong
				titleRaw = s.Find("td:nth-child(2) a").First().Text()
			}
			premiere.Title = pc.cleanTitle(titleRaw)
			if premiere.Title == "" {
				return //if there was no title then this is a garbage entry
			}
			if _, ok := premiereSet[premiere.Title]; ok {
				return //this is apparently a duplicate
			}

			//Check if it's a new series
			newImg := s.Find("td:nth-child(2) img[alt=\"new series\"]")
			if newImg != nil && newImg.Nodes != nil {
				premiere.IsNew = true
			}
			newImg = s.Find("td:nth-child(2) img[alt=\"limited series\"]")
			if newImg != nil && newImg.Nodes != nil {
				premiere.IsNew = true
			}

			//Get genres
			genresRaw, _ :=  s.Find("td:nth-child(2)").Html()
			parts := strings.Split(genresRaw, "<br/>")
			if len(parts) == 2 { 
				genresRaw = parts[1] //the genres are in the second part
				parts := strings.Split(genresRaw, ":")
				if len(parts) == 2 {
					genresRaw = parts[0]
				}
				genresRaw = strings.TrimSpace(genresRaw)
			}
			genreList := strings.Split(genresRaw, "/")
			if !pc.conf.IsInterestingMainGenre(genreList) {
				return //Not an interesting genre
			}
			premiere.Genres = genreList

			//Get streamer
			networkRaw := s.Find("td:nth-child(3)")
			premiere.StreamingOption = pc.getStreamer(networkRaw)

			premiereSet[premiere.Title] = premiere
		})
	})

	//turn the map into a list
	premieres := &PremiereList{
		StartDate: lastProcessedDate,
		EndDate:   newestDate,
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
	rawText := s.Text()
	if strings.Contains(rawText, "Netflix") {
		return streamer.Netflix
	}
	if strings.Contains(rawText, "Prime Video") {
		return streamer.Amazon
	}
	if strings.Contains(rawText, "Disney+") {
		return streamer.Disney
	}
	return streamer.None
}
