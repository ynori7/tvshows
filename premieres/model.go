package premieres

import "github.com/ynori7/tvshows/streamer"

type Premiere struct {
	Title           string
	IsNew           bool //if false, it's a new season of an older show
	Genres          []string
	StreamingOption streamer.Streamer
}

type PremiereList struct {
	StartDate string
	EndDate   string
	Premieres []Premiere
}
