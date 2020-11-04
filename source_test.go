package charityhonor

import (
	vinscraper "github.com/Vindexus/go-scraper"
	"testing"

	"github.com/monstercat/golib/expectm"
)

func TestParseSourceURL(t *testing.T) {
	type sourceTest struct {
		URL          string
		ExpectedType string
		ExpectedKey  string
		Error        error
		ExpectedMeta *expectm.ExpectedM
	}

	tests := []sourceTest{
		{
			URL:          "https://www.reddit.com/r/vancouver/comments/c78dd0/just_driving_the_wrong_way_on_a_highway_exit_with/",
			ExpectedType: vinscraper.SourceRedditPost,
			ExpectedKey:  "c78dd0",
			Error:        nil,
			ExpectedMeta: &expectm.ExpectedM{
				"subreddit": "vancouver",
				"author":    "shazoocow",
			},
		},
		{
			URL:          "https://np.reddit.com/r/pathofexile/comments/c6oy9e/to_everyone_that_feels_bored_by_the_game_or/esai27c/?context=3",
			ExpectedType: vinscraper.SourceRedditComment,
			ExpectedKey:  "esai27c",
			Error:        nil,
			ExpectedMeta: &expectm.ExpectedM{
				"subreddit": "pathofexile",
			},
		},
		{
			URL:          "https://www.reddit.com/about",
			Error:        nil,
			ExpectedType: vinscraper.SourceURL,
			ExpectedKey:  "https://www.reddit.com/about",
		},
		{
			URL:   "facebook colin",
			Error: vinscraper.ErrSourceInvalidURL,
		},
		{
			URL:   "twitter.com/@whatever",
			Error: vinscraper.ErrSourceInvalidURL,
		},
		{
			URL:          "http://twitter.com/@whatever",
			Error:        nil,
			ExpectedType: vinscraper.SourceURL,
			ExpectedKey:  "http://twitter.com/@whatever",
		},
	}

	for i, test := range tests {
		url := test.URL
		source, err := ParseSourceURL(url)
		if err != nil {
			if test.Error == nil {
				t.Fatal(err)
			} else {
				if test.Error != err {
					t.Errorf("#%d: Expected err %v but got %v", i, test.Error, err)
				}
				continue
			}
		}
		if source.Type != test.ExpectedType {
			t.Errorf("[%d] Type should be %v, found %v", i, test.ExpectedType, source.Type)
		}
		if source.Key != test.ExpectedKey {
			t.Errorf("[%d] Key should be %v, found %v", i, test.ExpectedKey, source.Key)
		}
		if test.ExpectedMeta != nil {
			if err != nil {
				t.Errorf("[%d] Error getting meta: %s", i, err)
			} else if err := expectm.CheckJSON(source.Meta, test.ExpectedMeta); err != nil {
				t.Errorf("[%d] %s", i, err)
			}
		}
	}
}
