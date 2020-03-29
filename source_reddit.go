package charityhonor

import (
	"errors"
	"net/url"
	"regexp"
	"strings"

	"github.com/monstercat/request"
)

type RedditPostInfo struct {
	Subreddit string `json:"subreddit"`
}

type RedditCommentInfo struct {
	Subreddit string `json:"subreddit"`
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

func (p *RedditPostSource) GetMeta() (interface{}, error) {
	return M{}, nil
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

func (p *RedditCommentSource) GetMeta() (interface{}, error) {
	var body RedditCommentInfoResponse
	params := request.Params{
		Url: "https://api.reddit.com/api/info?id=t1_" + p.ID,
	}
	if err := redditRequest(&params, nil, &body); err != nil {
		return nil, err
	}
	if len(body.Data.Children) == 0 {
		return nil, errors.New("no children in reddit")
	}
	dat := body.Data.Children[0].Data
	return dat, nil
}

func (p *RedditCommentSource) GetType() SourceType {
	return STRedditComment
}

func IsRedditSource(u *url.URL) bool {
	host := strings.ToLower(u.Hostname())
	if strings.HasSuffix(host, "reddit.com") {
		return true
	}
	return false
}

func ParseRedditSourceURL(urlS string) (Source, error) {
	r, err := regexp.Compile("\\/comments\\/([a-zA-Z0-9]+)\\/?[[a-zA-Z0-9\\_]+?\\/([a-zA-Z0-9]+)?")
	if err != nil {
		panic(err)
	}
	result := r.FindStringSubmatch(urlS)

	if len(result) < 3 {
		return NewDefaultSource(urlS), nil
	}

	if result[2] == "" {
		return &RedditPostSource{
			ID: result[1],
		}, nil
	}
	return &RedditCommentSource{
		ID: result[2],
	}, nil
}
