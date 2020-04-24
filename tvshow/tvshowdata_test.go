package tvshow

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ynori7/tvshows/config"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetTvShowData(t *testing.T) {
	//given
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		dat, err := ioutil.ReadFile("testdata/game-of-thrones.html")
		require.NoError(t, err, "There was an error reading the test data file")
		rw.Write(dat)
	}))
	defer server.Close()

	conf := config.Config{MainGenres: []string{"Drama", "Comedy"}}
	imdbClient := ImdbClient{httpClient: server.Client(), conf: conf, baseUrl: server.URL}

	//when
	tvShow, err := imdbClient.GetTvShowData(server.URL)

	//then
	require.NoError(t, err, "There was an error getting the link")
	assert.Equal(t, server.URL, tvShow.Link)
	assert.Equal(t, 5, len(tvShow.Genres))
	assert.Equal(t, "Game of Thrones", tvShow.Title)
	assert.Equal(t, 1663502, tvShow.Rating.RatingCount)
}

func Test_calculateScore(t *testing.T) {
	//given
	testcases := map[string]struct{
		Rating string
		RatingCount int
		Expected int
	} {
		"Game of Thrones, extremely well-rated": {
			Rating: "9.3",
			RatingCount: 1663502,
			Expected: 100,
		},
	}

	client := NewTvShowClient(config.Config{})

	for testcase, testdata := range testcases {
		//when
		score := client.calculateScore(testdata.Rating, testdata.RatingCount)

		//then
		assert.Equal(t, testdata.Expected, score, testcase)
	}
}