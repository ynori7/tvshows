package premieres

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ynori7/tvshows/config"
)

func Test_GetPotentiallyInterestingPremieres(t *testing.T) {
	//given
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		dat, err := ioutil.ReadFile("testdata/metacritic-tv-premieres.html")
		require.NoError(t, err, "There was an error reading the test data file")
		rw.Write(dat)
	}))
	defer server.Close()

	conf := config.Config{MainGenres: []string{"Drama", "Comedy"}}
	premieresClient := PremieresClient{httpClient: server.Client(), conf: conf, premieresUrl: server.URL}

	testcases := map[string]struct{
		date string
		expectedLen int
	} {
		"a week ago": {
			date:"June 4",
			expectedLen: 10,
		},
		"last date is in future": {
			date: "June 28",
			expectedLen: 90, //it'll process them all
		},
		"last date is same as most recent": {
			date: "June 11",
			expectedLen: 0,
		},
		"beginning of the month": {
			date: "June 1",
			expectedLen: 12,
		},
	}

	//when
	for testcase, testdata := range testcases {
		premieres, err := premieresClient.GetPotentiallyInterestingPremieres(testdata.date)

		//then
		require.NoError(t, err, "There was an error getting the premieres", testcase)
		assert.Equal(t, testdata.expectedLen, len(premieres.Premieres), testcase)
		assert.Equal(t, testdata.date, premieres.StartDate, testcase)
		assert.Equal(t, "June 11", premieres.EndDate, testcase)
	}
}

//Time bandits wrong genres
//Wandavision? what?