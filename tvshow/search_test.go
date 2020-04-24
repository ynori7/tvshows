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

func Test_Search(t *testing.T) {
	//given
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		dat, err := ioutil.ReadFile("testdata/search-got.html")
		require.NoError(t, err, "There was an error reading the test data file")
		rw.Write(dat)
	}))
	defer server.Close()

	conf := config.Config{MainGenres: []string{"Drama", "Comedy"}}
	imdbClient := NewTvShowClient(conf)
	imdbClient.httpClient = server.Client()
	imdbClient.baseUrl = server.URL

	//when
	link, err := imdbClient.SearchForTvSeriesTitle("Game of Thrones")

	//then
	require.NoError(t, err, "There was an error getting the link")
	assert.Equal(t, "https://www.imdb.com/title/tt0944947/", link)
}

func Test_Search_TitleNotExactMatch(t *testing.T) {
	//given
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		dat, err := ioutil.ReadFile("testdata/ghost-in-the-shell-sac_search.html")
		require.NoError(t, err, "There was an error reading the test data file")
		rw.Write(dat)
	}))
	defer server.Close()

	conf := config.Config{MainGenres: []string{"Drama", "Comedy"}}
	imdbClient := NewTvShowClient(conf)
	imdbClient.httpClient = server.Client()
	imdbClient.baseUrl = server.URL

	//when
	link, err := imdbClient.SearchForTvSeriesTitle("Ghost in the Shell: SAC_2045")

	//then
	require.NoError(t, err, "There was an error getting the link")
	assert.Equal(t, "https://www.imdb.com/title/tt9466298/", link)
}

func Test_Search_TitleHasAccent(t *testing.T) {
	//given
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		dat, err := ioutil.ReadFile("testdata/elite_search.html")
		require.NoError(t, err, "There was an error reading the test data file")
		rw.Write(dat)
	}))
	defer server.Close()

	conf := config.Config{MainGenres: []string{"Drama", "Comedy"}}
	imdbClient := NewTvShowClient(conf)
	imdbClient.httpClient = server.Client()
	imdbClient.baseUrl = server.URL

	//when
	link, err := imdbClient.SearchForTvSeriesTitle("Élite")

	//then
	require.NoError(t, err, "There was an error getting the link")
	assert.Equal(t, "https://www.imdb.com/title/tt7134908/", link)
}

func Test_fuzzifyTitle(t *testing.T) {
	testdata := map[string]string{
		"Ghost in the Shell: SAC_2045": "ghost in the shell sac2045",
		"Élite": "elite",
		"Manhunt: Deadly Games": "manhunt deadly games",
	}

	imdbClient := NewTvShowClient(config.Config{})

	for testcase, expected := range testdata {
		actual := imdbClient.fuzzifyTitle(testcase)
		assert.Equal(t, expected, actual)
	}
}