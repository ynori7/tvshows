package tvshow

import (
	"encoding/json"
	"github.com/ynori7/tvshows/streamer"
)

type TvShow struct {
	Title           string `json:"name"`
	Type            string `json:"@type"`
	Link            string
	Image           string `json:"image"`
	Genres          []string
	GenresRaw       interface{} `json:"genre"` //sometimes it's a string and sometimes it's a list of strings
	Rating          Rating      `json:"aggregateRating"`
	Description     string      `json:"description"`
	Created         string      `json:"datePublished"`
	AgeRating       string      `json:"contentRating"`
	Score           int
	StreamingOption streamer.Streamer
	IsNewSeries     bool
}

type Rating struct {
	AverageRating json.Number `json:"ratingValue"`
	RatingCount   int         `json:"ratingCount"`
}

type SearchResult struct {
	Title       string
	Link        string
	DedupNumber string //roman numeral to identify different shows with the same title
	Year        string
	Type        string //"Tv Series"
}
