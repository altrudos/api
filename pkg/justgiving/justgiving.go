package justgiving

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Mode string

var (
	ErrDonationNotFound = errors.New("Could not find donation")
)

var (
	ModeStaging    Mode = "staging"
	ModeProduction Mode = "production"
)

type Donation struct {
	Id                  int       `json:"int"`
	Amount              string    `json:"amount"`
	CurrencyCode        string    `json:"currencyCode"`
	LocalAmount         string    `json:"donorLocalAmount"`
	LocalCurrencyCode   string    `json:"donorLocalCurrencyCode"`
	Date                time.Time //This needs to be calculated by us cause they give us a weird date format
	DateString          string    `json:"donationDate"`
	ThirdPartyReference string    `json:"thirdPartyReference"`
	Status              string    `json:"status"`
}

type Charity struct {
	CountryCode string `json:"countryCode"`
	Description string `json:"description"`
	Id          int    `json:"id"`
	LogoUrl     string `db:"logo_url" json:"logo"`
	Name        string `json:"name"`
	Subtext     string `json:"subtext"`
	Summary     string `json:"summary"`
	WebsiteUrl  string `db:"website_url" json:"websiteUrl"`
}

type Params struct {
	Path   string
	Method string
	Url    string
	Body   string
	Debug  bool
	Query  url.Values
}

type JustGiving struct {
	Mode  Mode
	AppId string
	Debug bool
}

func GetStringInBetween(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str, end)
	return str[s:e]
}

func (d *Donation) GetDate() time.Time {
	num := GetStringInBetween(d.DateString, "Date(", "+")
	nanos, _ := strconv.Atoi(num)
	seconds := nanos / 1000
	date := time.Unix(int64(seconds), 0)
	return date
}

func (d *Donation) ConvertAmount(str string) float64 {
	amt, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return float64(0)
	}
	return amt

}

func (d *Donation) GetAmount() float64 {
	return d.ConvertAmount(d.Amount)
}

func (d *Donation) GetLocalAmount() float64 {
	return d.ConvertAmount(d.LocalAmount)
}

func (jg *JustGiving) Request(params *Params, send interface{}, receive interface{}) error {
	if jg.AppId == "" {
		return errors.New("No AppId found in JustGiving")
	}

	if jg.Mode == "" {
		return errors.New("No mode set. Must be either staging or production")
	}

	url := "https://api."

	if jg.Mode == ModeStaging {
		url += "staging."
	}

	url += "justgiving.com/" + jg.AppId + "/" + params.Path

	if len(params.Query) > 0 {
		url += "?" + params.Query.Encode()
	}

	req, err := http.NewRequest(params.Method, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	bodyString := string(bodyBytes)
	params.Body = bodyString

	if params.Debug || jg.Debug {
		fmt.Println("URL", url)
		fmt.Println("Body", bodyString)
	}

	err = json.Unmarshal(bodyBytes, &receive)
	if err != nil {
		return err
	}
	if res.StatusCode > 204 {
		err = errors.New(fmt.Sprintf("Status code too high: %d", res.StatusCode))
		return err
	}
	return nil
}

func (jg *JustGiving) GetCharityById(id int) (*Charity, error) {
	params := &Params{
		Path:   fmt.Sprintf("v1/charity/%d", id),
		Method: http.MethodGet,
	}

	charity := &Charity{}
	err := jg.Request(params, nil, charity)
	if err != nil {
		return nil, err
	}

	return charity, nil
}

func (jg *JustGiving) GetDonationByReference(reference string) (*Donation, error) {
	params := &Params{
		Path:   "v1/donation/ref/" + reference,
		Method: http.MethodGet,
	}

	type Response struct {
		Donations  []*Donation    `json:"donations"`
		Pagination map[string]int `json:"pagination"`
	}

	var resp Response
	err := jg.Request(params, nil, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Donations) == 0 {
		return nil, ErrDonationNotFound
	}

	if len(resp.Donations) > 1 {
		return nil, errors.New("Too many donations with that reference found.")
	}

	dono := resp.Donations[0]
	return dono, nil
}

func (jg *JustGiving) GetDonationById(id int) (*Donation, error) {
	params := &Params{
		Path:   "v1/donation/" + strconv.Itoa(id),
		Method: http.MethodGet,
	}

	var dono Donation
	err := jg.Request(params, nil, &dono)
	if err != nil {
		return nil, err
	}
	return &dono, nil
}

func (jg *JustGiving) GetDonationLink(charityId int, query url.Values) string {
	domain := "https://link."

	if jg.Mode == ModeStaging {
		domain += "staging."
	}

	domain += "justgiving.com"

	return fmt.Sprintf("%s/v1/charity/donate/charityId/%d?", domain, charityId) + query.Encode()
}
