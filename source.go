package altrudos

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

const (
	SourceTypeLink = "link"

	SourceTypeRedditPost = "reddit_post"

	SourceTypeYouTubeVideo   = "youtube_video"
	SourceTypeYouTubeChannel = "youtube_channel"
)

type Source struct {
	Type string
	Meta FlatMap
}

type DomainMatcher struct {
	Hosts []string
	Parse func(link *url.URL) (*Source, error)
}

func DefaultSource(link *url.URL) (*Source, error) {
	return &Source{
		Type: SourceTypeLink,
	}, nil
}

var DomainMatchers = []DomainMatcher{
	{
		Hosts: []string{"reddit.com", "redd.it"},
		Parse: func(link *url.URL) (*Source, error) {
			subreddit := regexp.MustCompile("\\/r\\/([a-zA-Z]+)\\/?")
			match := subreddit.FindAllString(link.String(), -1)

			// TODO: Check for comment

			if len(match) == 1 {
				return &Source{
					Type: SourceTypeRedditPost,
					Meta: FlatMap{
						"subreddit": match[0],
					},
				}, nil
			}

			return DefaultSource(link)
		},
	},
	{
		Hosts: []string{"youtube.com", "youtu.be"},
		Parse: func(link *url.URL) (*Source, error) {
			if link.Query().Get("v") != "" {
				return &Source{
					Type: SourceTypeYouTubeVideo,
				}, nil
			}

			if strings.Contains(link.Path, "/channel/") {
				return &Source{
					Type: SourceTypeYouTubeChannel,
				}, nil
			}

			return DefaultSource(link)
		},
	},
}

func ParseSourceURL(urlStr string) (*Source, error) {
	fmt.Println("urlStr", urlStr)
	parsed, err := url.Parse(urlStr)
	fmt.Println("parsed", parsed)
	fmt.Println("err", err)
	if err != nil {
		return nil, err
	}

	for _, dm := range DomainMatchers {
		found := false
		for _, domain := range dm.Hosts {
			if parsed.Host == domain {
				found = true
				break
			}
		}
		if !found {
			continue
		}

		return dm.Parse(parsed)
	}

	return &Source{
		Type: "link",
	}, nil
}
