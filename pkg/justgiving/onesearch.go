package justgiving

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

var (
	ErrGroupedResultsNot1 = errors.New("search did not return 1 grouped result")
)

type OneSearchResponse struct {
	GroupedResults []*OneSearchGroupedResult
}

type OneSearchGroupedResult struct {
	Title   string
	Count   int
	Results []*SearchResult
}

type SearchResult map[string]interface{}

type CharitySearchResponse struct {
	Title   string
	Count   int
	Results []*Charity
}

func ConvertSearchResultsToCharities(genericResults []*SearchResult) ([]*Charity, error) {
	j, err := json.Marshal(genericResults)
	if err != nil {
		return nil, err
	}

	var list []*Charity
	if err := json.Unmarshal(j, &list); err != nil {
		return nil, err
	}

	return list, nil
}

func (jg *JustGiving) SearchCharities(search string) (*CharitySearchResponse, error) {
	return jg.SearchCharitiesWithLimit(search, -1)
}

func (jg *JustGiving) SearchCharitiesWithLimit(search string, limit int) (*CharitySearchResponse, error) {
	qry := url.Values{
		"q": []string{search},
		"i": []string{"Charity"},
	}
	if limit > 0 {
		qry["limit"] = []string{strconv.Itoa(limit)}
	}

	params := &Params{
		Path:   "v1/onesearch",
		Method: http.MethodGet,
		Query:  qry,
		Debug:  jg.Debug,
	}

	var response OneSearchResponse
	err := jg.Request(params, nil, &response)
	if err != nil {
		return nil, err
	}

	if len(response.GroupedResults) != 1 {
		return nil, ErrGroupedResultsNot1
	}

	group := response.GroupedResults[0]

	charities, err := ConvertSearchResultsToCharities(group.Results)
	if err != nil {
		return nil, err
	}

	charityResult := &CharitySearchResponse{
		Count:   group.Count,
		Title:   group.Title,
		Results: charities,
	}

	return charityResult, nil
}
