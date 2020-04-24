package tvshow

import "github.com/ynori7/tvshows/streamer"

type TvShow struct {
	Title           string `json:"name"`
	Type            string `json:"@type"`
	Link            string
	Image           string   `json:"image"`
	Genres          []string `json:"genre"`
	Rating          Rating   `json:"aggregateRating"`
	Description     string   `json:"description"`
	Created         string   `json:"datePublished"`
	AgeRating       string   `json:"contentRating"`
	StreamingOption streamer.Streamer
}

type Rating struct {
	AverageRating string `json:"ratingValue"`
	RatingCount   int    `json:"ratingCount"`
}
