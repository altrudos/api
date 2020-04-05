package charityhonor

import (
	"errors"
	"net/url"
	"regexp"
	"strings"

	"github.com/monstercat/golib/request"
)

var (
	ErrRedditNoChildren = errors.New("no children in reddit")
)

var redditUrlRegexp = "\\/comments\\/([a-zA-Z0-9]+)\\/?[[a-zA-Z0-9\\_]+?\\/([a-zA-Z0-9]+)?"

type RedditPostInfo struct {
	RedditThing
	Title string `json:"title"`
}

type RedditCommentInfo struct {
	RedditThing
}

// Comment or Post
type RedditThing struct {
	Author    string  `json:"author"`
	Body      string  `json:"body"`
	Created   float64 `json:"created"`
	Permalink string  `json:"permalink"`
	Subreddit string  `json:"subreddit"`
}

func (c *RedditCommentInfo) ToMap() FlatMap {
	return FlatMap{
		"subreddit": c.Subreddit,
	}
}

func (c *RedditPostInfo) ToMap() FlatMap {
	return FlatMap{
		"subreddit": c.Subreddit,
	}
}

type RedditCommentInfoResponse struct {
	Data struct {
		Children []struct {
			Data RedditCommentInfo `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type RedditPostInfoResponse struct {
	Data struct {
		Children []struct {
			Data RedditPostInfo `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type RedditPostSource struct {
	ID string
}

type RedditCommentSource struct {
	ID string
}

func (p *RedditPostSource) GetKey() string {
	return p.ID
}

func (p *RedditPostSource) GetMeta() (FlatMap, error) {
	var body RedditPostInfoResponse
	params := request.Params{
		Url: "https://api.reddit.com/api/info?id=t3_" + p.ID,
	}
	if err := redditRequest(&params, nil, &body); err != nil {
		return nil, err
	}
	if len(body.Data.Children) == 0 {
		return nil, ErrRedditNoChildren
	}
	dat := body.Data.Children[0].Data
	return dat.ToMap(), nil
}

func (p *RedditPostSource) GetType() SourceType {
	return STRedditPost
}

func (p *RedditCommentSource) GetKey() string {
	return p.ID
}

func redditRequest(params *request.Params, payload interface{}, body interface{}) error {
	if params.Headers == nil {
		params.Headers = make(map[string]string)
	}
	params.Headers["User-agent"] = "charityhonor 0.1"
	return request.Request(params, payload, body)
}

func (p *RedditCommentSource) GetMeta() (FlatMap, error) {
	var body RedditCommentInfoResponse
	params := request.Params{
		Url: "https://api.reddit.com/api/info?id=t1_" + p.ID,
	}
	if err := redditRequest(&params, nil, &body); err != nil {
		return nil, err
	}
	if len(body.Data.Children) == 0 {
		return nil, ErrRedditNoChildren
	}
	dat := body.Data.Children[0].Data
	return dat.ToMap(), nil
}

func (p *RedditCommentSource) GetType() SourceType {
	return STRedditComment
}

func IsRedditSource(urlS string) bool {
	u, _ := url.Parse(urlS)
	host := strings.ToLower(u.Hostname())
	if !strings.HasSuffix(host, "reddit.com") {
		return false
	}
	r, err := regexp.Compile(redditUrlRegexp)
	if err != nil {
		return false
	}
	result := r.FindStringSubmatch(urlS)

	if len(result) < 3 {
		return false
	}
	return true
}

func ParseRedditSourceURL(urlS string) (Source, error) {
	r, err := regexp.Compile(redditUrlRegexp)
	if err != nil {
		return nil, err
	}
	result := r.FindStringSubmatch(urlS)

	if result[2] == "" {
		return &RedditPostSource{
			ID: result[1],
		}, nil
	}
	return &RedditCommentSource{
		ID: result[2],
	}, nil
}
