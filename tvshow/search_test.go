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
	imdbClient := ImdbClient{httpClient: server.Client(), conf: conf, baseUrl: server.URL}

	//when
	link, err := imdbClient.SearchForTvSeriesTitle("Game of Thrones")

	//then
	require.NoError(t, err, "There was an error getting the link")
	assert.Equal(t, "https://www.imdb.com/title/tt0944947/", link)
}
