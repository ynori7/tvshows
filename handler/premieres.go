package handler

import (
	"fmt"
	"github.com/ynori7/tvshows/config"
	"github.com/ynori7/tvshows/premieres"
	"io/ioutil"
	"strings"
	"time"
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

//todo return type
func (h PremieresHandler) GetNewPremieres() (*premieres.PremiereList, error) {
	lastProcessedDate := h.getLastProcessedDate()

	list, err := h.premieresClient.GetPotentiallyInterestingPremieres(lastProcessedDate)
	if err != nil {
		return nil, err
	}

	h.updateLastProcessedDate(list.EndDate)

	return list, nil
}

func (h PremieresHandler) getLastProcessedDate() string {
	dat, err := ioutil.ReadFile(lastProcessedFile)
	if err != nil || len(strings.TrimSpace(string(dat))) == 0 {
		lastWeek := time.Now().Add(-1 * defaultDays * 24 * time.Hour)
		return fmt.Sprintf("%s %d", lastWeek.Month(), lastWeek.Day())
	}

	return strings.TrimSpace(string(dat))
}

func (h PremieresHandler) updateLastProcessedDate(date string) error {
	return ioutil.WriteFile(lastProcessedFile, []byte(date), 0644)
}
