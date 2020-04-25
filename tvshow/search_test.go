package tvshow

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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
	imdbClient := NewImdbClient(conf)
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
	imdbClient := NewImdbClient(conf)
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
	imdbClient := NewImdbClient(conf)
	imdbClient.httpClient = server.Client()
	imdbClient.baseUrl = server.URL

	//when
	link, err := imdbClient.SearchForTvSeriesTitle("Élite")

	//then
	require.NoError(t, err, "There was an error getting the link")
	assert.Equal(t, "https://www.imdb.com/title/tt7134908/", link)
}

func Test_Search_TitleNonUnique(t *testing.T) {
	//given
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		dat, err := ioutil.ReadFile("testdata/sanctuary_search.html")
		require.NoError(t, err, "There was an error reading the test data file")
		rw.Write(dat)
	}))
	defer server.Close()

	conf := config.Config{MainGenres: []string{"Drama", "Comedy"}}
	imdbClient := NewImdbClient(conf)
	imdbClient.httpClient = server.Client()
	imdbClient.baseUrl = server.URL

	//when
	link, err := imdbClient.SearchForTvSeriesTitle("Sanctuary")

	//then
	require.NoError(t, err, "There was an error getting the link")
	assert.Equal(t, "https://www.imdb.com/title/tt8661868/", link)
}

func Test_Search_TitleNonUniqueAndBothSameYear(t *testing.T) {
	//given
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		dat, err := ioutil.ReadFile("testdata/the-stranger_search.html")
		require.NoError(t, err, "There was an error reading the test data file")
		rw.Write(dat)
	}))
	defer server.Close()

	conf := config.Config{MainGenres: []string{"Drama", "Comedy"}}
	imdbClient := NewImdbClient(conf)
	imdbClient.httpClient = server.Client()
	imdbClient.baseUrl = server.URL

	//when
	link, err := imdbClient.SearchForTvSeriesTitle("The Stranger")

	//then
	require.NoError(t, err, "There was an error getting the link")
	assert.Equal(t, "https://www.imdb.com/title/tt9698480/", link)
}

func Test_fuzzifyTitle(t *testing.T) {
	//given
	testdata := map[string]string{
		"Ghost in the Shell: SAC_2045": "ghost in the shell sac2045",
		"Élite":                        "elite",
		"Manhunt: Deadly Games":        "manhunt deadly games",
	}

	imdbClient := NewImdbClient(config.Config{})

	for testcase, expected := range testdata {
		//when
		actual := imdbClient.fuzzifyTitle(testcase)

		//then
		assert.Equal(t, expected, actual)
	}
}

func Test_parseSearchResultTitle(t *testing.T) {
	//given
	testdata := map[string]*SearchResult{
		"Sanctuary (2008) (TV Series)": {
			Title: "Sanctuary",
			Year:  "2008",
			Type:  "TV Series",
		},
		"Sanctuary (I) (2019) (TV Series)": {
			Title:       "Sanctuary",
			DedupNumber: "I",
			Year:        "2019",
			Type:        "TV Series",
		},
		"Money Heist (2017) (TV Series)": {
			Title: "Money Heist",
			Year:  "2017",
			Type:  "TV Series",
		},
		"The Terminator (1984)": nil, //not a series
		"The Stranger (II) (2020) (TV Series)": {
			Title:       "The Stranger",
			DedupNumber: "II",
			Year:        "2020",
			Type:        "TV Series",
		},
	}

	imdbClient := NewImdbClient(config.Config{})

	for testcase, expected := range testdata {
		//when
		searchParts := strings.Split(testcase, "(")
		searchTitle := searchParts[0]
		if expected != nil {
			searchTitle = expected.Title
		}
		actual := imdbClient.parseSearchResultTitle(strings.Split(testcase, "("), searchTitle)

		//then
		assert.Equal(t, expected, actual)
	}
}
