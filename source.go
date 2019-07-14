package charityhonor

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	ErrSourceInvalidURL    = errors.New("Invalid URL provided")
	ErrSourceInvalidReddit = errors.New("Unrecognized reddit link format")
)

type SourceType string
type SourceID interface{}

var (
	STRedditPost    SourceType = "reddit_post"
	STRedditComment SourceType = "reddit_comment"
	STURL           SourceType = "url" //For basically anything we don't know
)

type Source struct {
	Type SourceType
	ID   SourceID
}

func (s *Source) String() string {
	return fmt.Sprintf("%v:%s", s.Type, s.ID)
}

/**
 * A URL for some content to be honored is turned into a "Source"
 * This normalizes different URLs whose structure we know points
 * to the same content.
 * youtu.be/1234 and youtube.com/watch?v=1234 => youtube:1234 (not yet implemented)
 */
func ParseSourceURL(urlStr string) (*Source, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, ErrSourceInvalidURL
	}

	host := u.Hostname()

	if host == "" {
		return nil, ErrSourceInvalidURL
	}

	host = strings.ToLower(host)

	if strings.HasSuffix(host, "reddit.com") {
		return ParseRedditSourceURL(urlStr)
	}

	return &Source{
		Type: STURL,
		ID:   urlStr,
	}, nil
}

func ParseRedditSourceURL(urlStr string) (*Source, error) {
	//																post id      /post_title_here/comment_id
	r, err := regexp.Compile("\\/comments\\/([a-zA-Z0-9]+)\\/?[[a-zA-Z0-9\\_]+?\\/([a-zA-Z0-9]+)?")
	if err != nil {
		panic(err)
	}
	result := r.FindStringSubmatch(urlStr)

	if len(result) < 3 {
		return &Source{
			Type: STURL,
			ID:   urlStr,
		}, nil
	}

	if result[2] == "" {
		return &Source{
			Type: STRedditPost,
			ID:   result[1],
		}, nil
	} else {
		return &Source{
			Type: STRedditComment,
			ID:   result[2],
		}, nil
	}

	return &Source{
		Type: STRedditPost,
	}, nil
}
