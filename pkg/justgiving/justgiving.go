package justgiving

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var mode = "sandbox"

type DonationDetails struct {
}

type Params struct {
	Path   string
	Method string
}

type JustGiving struct {
	Mode  string
	AppId string
}

func (jg *JustGiving) Request(params *Params, data interface{}, body interface{}) error {
	url := "https://api."

	if jg.Mode == "sandbox" {
		url += "sandbox."
	}

	url += ".justgiving.com/" + jg.AppId + "/" + params.Path

	req, err := http.NewRequest(params.Method, params.Url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return err
	}
	if res.StatusCode > 204 {
		err = errors.New("Status code too high")
		return err
	}
	return nil
}

func (jg *JustGiving) RetrieveDonationDetailsByReference(reference string) {
	params := Params{
		Path:   "/v1/donation/ref/" + reference,
		Method: http.MethodGet,
	}

	donation := Get
}

func (jg *JustGiving) RetrieveDonationDetails(reference string) {
	params := Params{
		Path:   "/v1/donation/ref/" + reference,
		Method: http.MethodGet,
	}

	donation := Get
}
