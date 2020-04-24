package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	testConfig := []byte(`title: "rap-and-metal"
main_genres: #these genres are filtered first from the new releases list
  - "Drama"
  - "Comedy"
  - "Thriller"
  - "Animation"
  - "Family"
  - "Horror"
  - "Fantasy"
  - "Action"
  - "Sci-fi"
  - "Anime"
sub_genres: #these genres are filtered on the detail pages
  fuzzy_matches:
    - "Drama"
    - "Horror"
  exact_matches:
    - "Anime"
email:
  enabled: true
  private_key: "private123"
  public_key: "public456"
  from:
    address: "no-reply@something.com"
    name: "Nobody"
  to:
    address: "me@mysite.com"
    name: "Me"`)

	c := Config{}

	err := c.Parse(testConfig)
	require.NoError(t, err, "It should parse the config successfully")

	assert.Equal(t, "rap-and-metal", c.Title)
	assert.Equal(t, 10, len(c.MainGenres))
	assert.Equal(t, 2, len(c.SubGenres.FuzzyMatches))
	assert.Equal(t, 1, len(c.SubGenres.ExactMatches))
	assert.True(t, c.Email.Enabled)
	assert.Equal(t, c.Email.PrivateKey, "private123")
	assert.Equal(t, c.Email.PublicKey, "public456")
	assert.Equal(t, c.Email.From.Address, "no-reply@something.com")
	assert.Equal(t, c.Email.To.Name, "Me")
}

func Test_IsInterestingMainGenre(t *testing.T) {
	testcases := map[string]struct {
		List     []string
		Genres   []string
		Expected bool
	}{
		"No match": {
			List:     []string{"Drama", "Comedy", "Horror"},
			Genres:   []string{"Reality"},
			Expected: false,
		},
		"Exact match": {
			List:     []string{"Drama", "Comedy", "Horror"},
			Genres:   []string{"Drama"},
			Expected: true,
		},
		"Fuzzy match": {
			List:     []string{"Drama", "Comedy", "Horror"},
			Genres:   []string{"Comedy special"},
			Expected: false,
		},
		"Empty list": {
			List:     []string{},
			Genres:   []string{"Drama"},
			Expected: false,
		},
	}

	for testcase, testdata := range testcases {
		c := Config{MainGenres: testdata.List}
		res := c.IsInterestingMainGenre(testdata.Genres)
		assert.Equal(t, testdata.Expected, res, testcase)
	}
}

func Test_IsInterestingSubGenre(t *testing.T) {
	testcases := map[string]struct {
		FuzzyList []string
		ExactList []string
		Genre     string
		Expected  bool
	}{
		"No match": {
			FuzzyList: []string{"Jazz", "Country", "Pop/Rock"},
			ExactList: []string{"Grunge"},
			Genre:     "Rap",
			Expected:  false,
		},
		"Exact match in fuzzy list": {
			FuzzyList: []string{"Jazz", "Rap", "Rock"},
			ExactList: []string{"Grunge"},
			Genre:     "Rap",
			Expected:  true,
		},
		"Exact match in exact list": {
			FuzzyList: []string{"Jazz", "Rap", "Rock"},
			ExactList: []string{"Grunge"},
			Genre:     "Grunge",
			Expected:  true,
		},
		"Fuzzy match": {
			FuzzyList: []string{"Jazz", "Rap", "Rock"},
			ExactList: []string{"Grunge"},
			Genre:     "Pop/Rock",
			Expected:  true,
		},
		"Fuzzy match of exact list": {
			FuzzyList: []string{"Jazz", "Rap", "Rock"},
			ExactList: []string{"Grunge"},
			Genre:     "Post-Grunge",
			Expected:  false,
		},
		"Empty list": {
			FuzzyList: []string{},
			ExactList: []string{},
			Genre:     "Rock",
			Expected:  false,
		},
	}

	for testcase, testdata := range testcases {
		c := Config{}
		c.SubGenres.FuzzyMatches = testdata.FuzzyList
		c.SubGenres.ExactMatches = testdata.ExactList

		res := c.IsInterestingSubGenre(testdata.Genre)
		assert.Equal(t, testdata.Expected, res, testcase)
	}
}
