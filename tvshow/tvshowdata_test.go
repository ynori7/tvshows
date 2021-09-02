package tvshow

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ynori7/tvshows/config"
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
	assert.Equal(t, 3, len(tvShow.Genres))
	assert.Equal(t, "Game of Thrones", tvShow.Title)
	assert.Equal(t, 1865597, tvShow.Rating.RatingCount)
}

func Test_GetTvShowData_GenresNotList(t *testing.T) {
	//given
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		dat, err := ioutil.ReadFile("testdata/blackaf.html")
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
	assert.Equal(t, 1, len(tvShow.Genres))
	assert.Equal(t, "#BlackAF", tvShow.Title)
	assert.Equal(t, 1516, tvShow.Rating.RatingCount)
}

func Test_calculateScore(t *testing.T) {
	//given
	testcases := map[string]struct {
		Rating      string
		RatingCount int
		Expected    int
	}{
		"Game of Thrones, extremely well-rated": {
			Rating:      "9.3",
			RatingCount: 1663502,
			Expected:    100,
		},
		"#BlackAF, new show with not very high rating": {
			Rating:      "6.6",
			RatingCount: 1516,
			Expected:    39,
		},
		"Tokyo Ghoul, decent rating with a fairly high count": {
			Rating:      "7.9",
			RatingCount: 30248,
			Expected:    79,
		},
		"Food Wars, good rating, but not so many": {
			Rating:      "8.2",
			RatingCount: 4279,
			Expected:    69,
		},
		"The Bachelor, terrible rating, middle count": {
			Rating:      "3.2",
			RatingCount: 6026,
			Expected:    27,
		},
	}

	client := NewImdbClient(config.Config{})

	for testcase, testdata := range testcases {
		//when
		score := client.calculateScore(testdata.Rating, testdata.RatingCount)

		//then
		assert.Equal(t, testdata.Expected, score, testcase)
	}
}
