package charityhonor

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"

	dbUtil "github.com/monstercat/golib/db"

	"github.com/monstercat/pgnull"

	"github.com/charityhonor/ch-api/pkg/justgiving"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

var (
	ErrMissingReferenceCode = errors.New("donation is missing reference code")
	ErrInvalidAmount = errors.New("invalid donation amount")
	ErrNegativeAmount = errors.New("donation amount can't be negative")
)

var (
	TABLE_DONATIONS = "donations"
)

var (
	// These statuses are uppercased on JustGiving s well
	DonationAccepted DonationStatus = "Accepted"
	DonationPending  DonationStatus = "Pending"
	DonationRejected DonationStatus = "Rejected"
)

var (
	ErrNoCurrencyCode = errors.New("no currency code")
	ErrNoAmount       = errors.New("no donation amount")
	ErrNoCharity      = errors.New("no charity")
)

var (
	DonationInsertBuilder = QueryBuilder.Insert(TABLE_DONATIONS)
	DonationUpdateBuilder = QueryBuilder.Update(TABLE_DONATIONS)
)

var (
	DONATION_COLUMNS = map[string]string{
		"Id":                "id",
		"CharityId":         "charity_id",
		"Created":           "created",
		"DonorAmount":       "donor_amount",
		"DonorCurrentyCode": "donor_currency_code",
		"DonorName":         "donor_name",
		"DriveId":           "drive_id",
		"FinalAmount":       "final_amount",
		"LastChecked":       "last_checked",
		"Message":           "message",
		"ReferenceCode":     "reference_code",
		"Status":            "status",
	}
)

var codeCount = 1

/*
"amount": "2.00",
    "currencyCode": "GBP",
    "donationDate": "\/Date(1556326412351+0000)\/",
    "donationRef": null,
    "donorDisplayName": "Awesome Guy",
    "donorLocalAmount": "2.75",
    "donorLocalCurrencyCode": "EUR",
    "donorRealName": "Peter Queue",
    "estimatedTaxReclaim": 0.56,
    "id": 1234,
    "image": "",
    "message": "Hope you like my donation. Rock on!",
    "source": "SponsorshipDonations",
    "status": "Accepted",
    "thirdPartyReference": "1234-my-sdi-ref"
*/
type DonationStatus string

type Donation struct {
	Charity           *Charity `db:"-"`
	CharityId         string   `db:"charity_id"`
	Created           time.Time
	DonorAmount       int               `db:"donor_amount"`        // What the donor typed in
	DonorCurrencyCode string            `db:"donor_currency_code"` // What the donor selected
	DonorName         pgnull.NullString `db:"donor_name"`
	DriveId           string            `db:"drive_id"`
	FinalAmount       int               `db:"final_amount"`
	Id                string            `setmap:"omitinsert"`
	LastChecked       pgnull.NullTime   `db:"last_checked"`
	Message           pgnull.NullString `db:"message"`
	Status            DonationStatus
	ReferenceCode     string `db:"reference_code"`
}

// Used in queries
type DonationOperators struct {
	*BaseOperator
	Statuses []DonationStatus
}

func GetDonationByField(tx sqlx.Queryer, field string, val interface{}) (*Donation, error) {
	query, args, err := QueryBuilder.
		Select(GetColumns(DONATION_COLUMNS)...).
		From(TABLE_DONATIONS).Where(field+"=?", val).
		ToSql()
	if err != nil {
		return nil, err
	}

	var d Donation
	err = sqlx.Get(tx, &d, query, args...)
	if err != nil {
		return nil, err
	}

	if d.CharityId == "" {
		return nil, errors.New("charity has a blank ID")
	}

	charity, err := GetCharityById(tx, d.CharityId)
	if err != nil {
		return nil, err
	}

	d.Charity = charity
	return &d, nil
}

func GetDonationById(tx sqlx.Queryer, id string) (*Donation, error) {
	return GetDonationByField(tx, "id", id)
}

func GetDonationByReferenceCode(tx sqlx.Queryer, code string) (*Donation, error) {
	return GetDonationByField(tx, "reference_code", code)
}

func GetDonationsToCheck(tx sqlx.Queryer) ([]*Donation, error) {
	ops := &DonationOperators{
		BaseOperator: &BaseOperator{
			SortField: "next_check",
			SortDir:   SortAsc,
		},
		Statuses: []DonationStatus{DonationPending},
	}

	return GetDonations(tx, ops)
}

func QueryDonations(q sqlx.Queryer, query *squirrel.SelectBuilder) ([]*Donation, error) {
	s, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	donos := make([]*Donation, 0)
	err = sqlx.Select(q, &donos, s, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return donos, nil
		}
		return nil, err
	}
	return donos, nil
}

func GetDonations(q sqlx.Queryer, ops *DonationOperators) ([]*Donation, error) {
	query := QueryBuilder.
		Select(GetColumns(DONATION_COLUMNS)...).
		From(TABLE_DONATIONS)

	if len(ops.Statuses) > 0 {
		query = query.Where("status = ANY (?)", StatusesPQStringArray(ops.Statuses))
	}

	return QueryDonations(q, &query)
}

func (d *Donation) GenerateReferenceCode(ext sqlx.Ext) error {
	exists := false
	for d.ReferenceCode == "" || exists == true {
		str := uuid.Must(uuid.NewV4()).String()
		str = fmt.Sprintf("ch-%d", time.Now().UnixNano())
		d.ReferenceCode = str
		dupe, err := GetDonationByReferenceCode(ext, d.ReferenceCode)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			}
			return err
		}
		exists = dupe != nil
	}

	return nil
}

//Create does magic before insert into db
func (d *Donation) Create(ext sqlx.Ext) error {
	if d.CharityId == "" {
		return ErrNoCharity
	}
	charity, err := GetCharityById(ext, d.CharityId)
	if err != nil {
		fmt.Println("err", err)
		if err == sql.ErrNoRows {
			return ErrCharityNotFound
		}
		return err
	}

	d.Charity = charity
	if d.ReferenceCode != "" {
		return ErrAlreadyInserted
	}

	if err := d.Validate(); err != nil {
		return err
	}

	err = d.GenerateReferenceCode(ext)
	if err != nil {
		return err
	}

	if d.Status == DonationStatus("") {
		d.Status = DonationPending
	}

	return d.Insert(ext)
}

func (d *Donation) Validate() error {
	currency, err := ParseCurrency(d.DonorCurrencyCode)
	if err != nil {
		return err
	}
	d.DonorCurrencyCode = currency

	if d.DonorAmount < 0 {
		return ErrNegativeAmount
	}

	return nil
}

//Raw insert into db
func (d *Donation) Insert(ext sqlx.Ext) error {
	query := DonationInsertBuilder.
		SetMap(dbUtil.SetMap(d, true)).
		Suffix(RETURNING_ID)
	return query.
		RunWith(ext).
		QueryRow().
		Scan(&d.Id)
}

func (d *Donation) Save(ext sqlx.Ext) error {
	setMap := dbUtil.SetMap(d, false)
	_, err := DonationUpdateBuilder.
		SetMap(setMap).
		Where("id=?", d.Id).
		RunWith(ext).
		Exec()
	return err
}

/*https://link.justgiving.com/v1/charity/donate/charityId/2096
?amount=10.00
&currency=USD
&reference=89302483&
exitUrl=http%3A%2F%2Flocalhost%3A9000%2Fconfirm%2F8930248302840%3FjgDonationId%3DJUSTGIVING-DONATION-ID
&message=Woohoo!%20Let's%20fight%20cancer!
*/
func (d *Donation) GetDonationLink(jg *justgiving.JustGiving) (string, error) {
	urls := url.Values{}
	if d == nil {
		panic("donation is nil")
	}
	fmt.Println("message", d.Message)
	if d.Message.Valid && d.Message.String != "" {
		urls.Set("message", d.Message.String)
	}

	if d.Message.Valid && d.Message.String != "" {
		urls.Set("message", d.Message.String)
	}

	if d.DonorCurrencyCode == "" {
		return "", ErrNoCurrencyCode
	}

	if d.DonorAmount == 0 {
		return "", ErrNoAmount
	}

	if d.Charity == nil {
		return "", ErrNoCharity
	}

	urls.Set("currency", d.DonorCurrencyCode)
	urls.Set("amount", AmountToString(d.DonorAmount))
	urls.Set("reference", d.ReferenceCode)

	return jg.GetDonationLink(d.Charity.JustGivingCharityId, urls), nil
}

func (d *Donation) GetJustGivingDonation(jg *justgiving.JustGiving) (*justgiving.Donation, error) {
	return jg.GetDonationByReference(d.ReferenceCode)
}

func (d *Donation) GetLastChecked() time.Time {
	if d.LastChecked.Valid {
		val, err := d.LastChecked.Value()
		if err != nil {
			return time.Time{}
		}

		return val.(time.Time)
	}

	return time.Time{}
}

func (d *Donation) AmountString() string {
	return AmountToString(d.FinalAmount)
}

func (d *Donation) IsAnonymous() bool {
	return !d.DonorName.Valid || d.DonorName.String == ""
}

func (d *Donation) GetDonorName() string {
	if !d.IsAnonymous() {
		return d.DonorName.String
	}
	return "Anonymous"
}

func (d *Donation) CheckStatus(tx *sqlx.Tx, jg *justgiving.JustGiving) error {
	jgDonation, err := d.GetJustGivingDonation(jg)
	var status DonationStatus
	if err != nil {
		if err == justgiving.ErrDonationNotFound {
			status = DonationRejected
		} else {
			return err
		}
	} else {
		status = DonationStatus(jgDonation.Status)
	}

	d.Status = status
	d.LastChecked = pgnull.NullTime{time.Now(), true}
	err = d.Save(tx)
	return err
}

func ApplyApproved(q *squirrel.SelectBuilder) {
	*q = q.Where("status=?", DonationAccepted)
}

// Takes an amount that is string from the frontend and returns it in cents
func AmountFromString(amount string) (int, error) {
	f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, ErrInvalidAmount
	}

	if f < 0 {
		return 0, ErrNegativeAmount
	}

	// Convert dollars to cents
	return int(f * 100), nil
}