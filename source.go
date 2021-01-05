package altrudos

import (
	"encoding/json"
	"errors"

	vinscraper "github.com/Vindexus/go-scraper"
)

var (
	ErrSourceInvalidReddit = errors.New("Unrecognized reddit link format")
)

type Source struct {
	Type string
	Key  string
	Meta FlatMap
}

func NewScraper() *vinscraper.Scraping {
	return &vinscraper.Scraping{
		// The order matters
		Scrapers: []vinscraper.Scraper{
			&vinscraper.RedditScraper{
				UserAgent: "altrudos-1.0",
			},
			&vinscraper.ScraperGeneric{},
		},
		TitleReplacers: []vinscraper.ScrapeReplacer{},
	}
}

func ParseSourceURL(urlStr string) (*Source, error) {
	scraper := NewScraper()

	info, err := scraper.Scrape(urlStr)
	if err != nil {
		return nil, err
	}

	// Conver the interface into a regular flatmap
	meta := FlatMap{}
	bytes, err := json.Marshal(info.Meta)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &meta); err != nil {
		return nil, err
	}

	// Some of the standard stuff the scraper returns is part
	// of the generic ScrapeInfo object, and not necessarily stored
	// in the meta object
	// For example a reddit post has the title in info.Title, and not in meta["Title"]
	// we need to copy this over because we only deal with the Meta in Altrudos
	if _, ok := meta["Title"]; !ok {
		meta["Title"] = info.Title
	}

	s := Source{
		Type: string(info.SourceType),
		Key:  info.SourceKey,
		Meta: meta,
	}

	return &s, nil
}
