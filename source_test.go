package charityhonor

import "testing"

func TestParseSourceURL(t *testing.T) {
	type sourceTest struct {
		URL          string
		ExpectedType SourceType
		ExpectedID   SourceID
		Error        error
	}

	tests := []sourceTest{
		{
			URL:          "https://www.reddit.com/r/vancouver/comments/c78dd0/just_driving_the_wrong_way_on_a_highway_exit_with/",
			ExpectedType: STRedditPost,
			ExpectedID:   "c78dd0",
			Error:        nil,
		},
		{
			URL:          "https://np.reddit.com/r/pathofexile/comments/c6oy9e/to_everyone_that_feels_bored_by_the_game_or/esai27c/?context=3",
			ExpectedType: STRedditComment,
			ExpectedID:   "esai27c",
			Error:        nil,
		},
		{
			URL:          "https://www.reddit.com/about",
			Error:        nil,
			ExpectedType: STURL,
			ExpectedID:   "https://www.reddit.com/about",
		},
		{
			URL:   "facebook colin",
			Error: ErrSourceInvalidURL,
		},
		{
			URL:   "twitter.com/@whatever",
			Error: ErrSourceInvalidURL,
		},
		{
			URL:          "http://twitter.com/@whatever",
			Error:        nil,
			ExpectedType: STURL,
			ExpectedID:   "http://twitter.com/@whatever",
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
			t.Errorf("#%d: Type should be %v, found %v", i, test.ExpectedType, source.Type)
		}
		if source.ID != test.ExpectedID {
			t.Errorf("#%d: ID should be %v, found %v", i, test.ExpectedID, source.ID)
		}
	}
}
