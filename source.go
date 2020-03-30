package charityhonor

import (
	"errors"
	"net/url"
)

var (
	ErrSourceInvalidURL    = errors.New("Invalid URL provided")
	ErrSourceInvalidReddit = errors.New("Unrecognized reddit link format")
)

type SourceType string

var (
	STRedditPost    SourceType = "reddit_post"
	STRedditComment SourceType = "reddit_comment"
	STURL           SourceType = "url" //For basically anything we don't know
)

type Source interface {
	GetType() SourceType
	GetKey() string
	GetMeta() (FlatMap, error)
}

type DefaultSource struct {
	Type SourceType
	URL  string
	Meta M
}

func (s *DefaultSource) String() string {
	return s.URL
}

func (s *DefaultSource) GetType() SourceType {
	return STURL
}

func (s *DefaultSource) GetKey() string {
	return s.URL
}

func (s *DefaultSource) GetMeta() (FlatMap, error) {
	// TODO: Fetch the page's HTML and look for page title and maybe og: tags
	return FlatMap{
		"url": s.URL,
	}, nil
}

func NewDefaultSource(url string) *DefaultSource {
	return &DefaultSource{
		Type: STURL,
		URL:  url,
	}
}

/**
* A URL for some content to be honored is turne			d into a "Source"
 * This normalizes different URLs whose structure we know points
 * to the same content.
 * youtu.be/1234 and youtube.com/watch?v=1234 => youtube:1234 (not yet implemented)
*/
func ParseSourceURL(urlStr string) (Source, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, ErrSourceInvalidURL
	}

	host := u.Hostname()
	if host == "" {
		return nil, ErrSourceInvalidURL
	}

	if IsRedditSource(urlStr) {
		return ParseRedditSourceURL(urlStr)
	}

	return NewDefaultSource(urlStr), nil
}
