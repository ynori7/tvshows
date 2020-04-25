package tvshow

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/ynori7/tvshows/config"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const (
	baseUrl   = "https://www.imdb.com"
	searchURI = "/find"
)

type ImdbClient struct {
	httpClient *http.Client
	conf       config.Config
	baseUrl    string
	titleRegex *regexp.Regexp
}

func NewImdbClient(conf config.Config) ImdbClient {
	reg, err := regexp.Compile("[^a-zA-Z0-9\\s]+")
	if err != nil {
		log.Fatal(err)
	}

	return ImdbClient{
		httpClient: &http.Client{},
		conf:       conf,
		baseUrl:    baseUrl,
		titleRegex: reg,
	}
}

//GetTvShowData looks up the tv show details
func (c ImdbClient) GetTvShowData(link string) (*TvShow, error) {
	// Request the HTML page.
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Language", "en-US")
	res, err := c.httpClient.Do(req)
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

	jsonRaw := doc.Find("script[type=\"application/ld+json\"]")
	x := jsonRaw.Text()

	tvShow := new(TvShow)
	if err := json.Unmarshal([]byte(x), tvShow); err != nil {
		return nil, err
	}
	if tvShow.GenresRaw != nil {
		genres, ok := tvShow.GenresRaw.([]interface{})
		if ok {
			for _, g := range genres {
				tvShow.Genres = append(tvShow.Genres, g.(string))
			}
		} else {
			tvShow.Genres = []string{tvShow.GenresRaw.(string)}
		}
	}

	tvShow.Link = link
	tvShow.Score = c.calculateScore(tvShow.Rating.AverageRating, tvShow.Rating.RatingCount)

	return tvShow, nil
}

//SearchForTvSeriesTitle returns the IMDB url for the title
func (c ImdbClient) SearchForTvSeriesTitle(searchTitle string) (string, error) {
	// Request the HTML page.
	req, err := http.NewRequest("GET", c.buildImdbSearchUrl(searchTitle), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept-Language", "en-US")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	foundLink := ""
	found := false
	searchTitle = c.fuzzifyTitle(searchTitle)

	// Find the new releases
	doc.Find("table.findList tr").Each(func(i int, s *goquery.Selection) {
		if found {
			return
		}

		res := s.Find(".result_text")
		resText := res.Text()
		textParts := strings.Split(resText, "(")
		title := strings.TrimSpace(textParts[0])
		if c.fuzzifyTitle(title) != searchTitle {
			return //not an exact match
		}
		if len(textParts) != 3 {
			return //it's not a tv show
		}

		resType := strings.Trim(textParts[2], ") ")
		if !strings.HasPrefix(resType, "TV Series") && !strings.HasPrefix(resType, "TV Mini-Series") {
			return //not a tv show
		}

		linkRaw := res.Find("a")
		if link, ok := linkRaw.Attr("href"); ok {
			found = true
			foundLink = baseUrl + link
		}
	})

	if foundLink == "" {
		return "", fmt.Errorf("no result found")
	}

	return foundLink, nil
}

//fuzzifyTitle normalizes the text by removing punctuation and accents to make the titles comparable
func (c ImdbClient) fuzzifyTitle(t string) string {
	//replace accented characters
	tr := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
	}), norm.NFC)
	result, _, _ := transform.String(tr, t)

	result = c.titleRegex.ReplaceAllString(result, "") //remove punctuation

	return strings.ToLower(result)
}

func (c ImdbClient) buildImdbSearchUrl(title string) string {
	params := url.Values{}
	params.Add("q", title)
	return fmt.Sprintf("%s%s?%s", c.baseUrl, searchURI, params.Encode())
}

//list of rating counts. The index is the log() value
var scoreIntervals = []int{
	0, 0, 500, 1000, 1500, 2000, 3000, 4000, 8000, 10000, 20000, 50000, 100000, 500000,
}

//calculateScore returns a score out of 100 based on the imdb rating, weight by the number of ratings
func (c ImdbClient) calculateScore(averageRating string, ratingCount int) int {
	rating, err := strconv.ParseFloat(averageRating, 64)
	if err != nil {
		return 0
	}

	scoreBase := rating * float64(10)
	score := 0

	for i := len(scoreIntervals) - 1; i > 0; i-- {
		if ratingCount >= scoreIntervals[i] {
			score = int(scoreBase * math.Log10(float64(i)))
			break
		}
	}

	if score > 100 {
		return 100
	}
	return score
}
