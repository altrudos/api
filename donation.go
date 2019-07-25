package charityhonor

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/satori/go.uuid"
	"github.com/charityhonor/ch-api/pkg/justgiving"
)

var (
	ErrMissingReferenceCode = errors.New("Donation is missing reference code")
)

var (
	TABLE_DONATIONS = "donations"
)

var (
	DonationAccepted DonationStatus = "Accepted"
	DonationPending  DonationStatus = "Pending"
	DonationRejected DonationStatus = "Rejected"
)

var (
	DonationInsertBuilder = QueryBuilder.Insert(TABLE_DONATIONS)
	DonationUpdateBuilder = QueryBuilder.Update(TABLE_DONATIONS)
)

var (
	DONATION_COLUMNS = map[string]string{
		"Id":            "id",
		"LastChecked":   "last_checked",
		"Status":        "status",
		"ReferenceCode": "reference_code",
		"DriveId":       "drive_id",
		"Created":       "created",
		"DonorName":     "donor_name",
		"CharityId":     "charity_id",
		"Amount":        "amount",
		"CurrentyCode":  "currency_code",
		"Message":       "message",
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
	Id            int
	Created       time.Time
	LastChecked   pq.NullTime `db:"last_checked"`
	Status        DonationStatus
	ReferenceCode string  `db:"reference_code"`
	DriveId       int     `db:"drive_id"`
	CharityId     int     `db:"charity_id"`
	Amount        float64 `db:"amount"`
	DonorName     string  `db:"donor_name"`
	Message       string  `db:"message"`
	CurrencyCode  string  `db:"currency_code"`
	Charity       *Charity
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

	if d.CharityId == 0 {
		return nil, errors.New("Charity has an ID of 0")
	}

	charity, err := GetCharityById(tx, d.CharityId)
	if err != nil {
		return nil, err
	}

	d.Charity = charity
	return &d, nil
}

func GetDonationById(tx sqlx.Queryer, id int) (*Donation, error) {
	return GetDonationByField(tx, "id", id)
}

func GetDonationByReferenceCode(tx sqlx.Queryer, code string) (*Donation, error) {
	return GetDonationByField(tx, "reference_code", code)
}

func (d *Donation) GenerateReferenceCode(tx *sqlx.Tx) error {
	exists := false
	for d.ReferenceCode == "" || exists == true {
		str := uuid.Must(uuid.NewV4()).String()
		str = fmt.Sprintf("ch-%d", time.Now().UnixNano())
		d.ReferenceCode = str
		dupe, err := GetDonationByReferenceCode(tx, d.ReferenceCode)
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

func (d *Donation) getSetMap() M {
	return M{
		"charity_id":     d.CharityId,
		"drive_id":       d.DriveId,
		"status":         d.Status,
		"donor_name":     d.DonorName,
		"message":        d.Message,
		"amount":         d.Amount,
		"reference_code": d.ReferenceCode,
		"currency_code":  d.CurrencyCode,
	}
}

//Create does magic before insert into db
func (d *Donation) Create(tx *sqlx.Tx) error {
	charity, err := GetCharityById(tx, d.CharityId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrCharityNotFound
		}
		return err
	}

	d.Charity = charity

	if d.ReferenceCode != "" {
		return ErrAlreadyInserted
	}

	err = d.GenerateReferenceCode(tx)

	if err != nil {
		return err
	}

	return d.Insert(tx)
}

//Raw insert into db
func (d *Donation) Insert(tx *sqlx.Tx) error {
	setMap := d.getSetMap()
	return DonationInsertBuilder.
		SetMap(setMap).
		Suffix(RETURNING_ID).
		RunWith(tx).
		QueryRow().
		Scan(&d.Id)
}

func (d *Donation) Save(tx *sqlx.Tx) error {
	setMap := d.getSetMap()
	fmt.Println("setMap", setMap)
	fmt.Println("d.Id", d.Id)
	_, err := DonationUpdateBuilder.
		SetMap(setMap).
		Where("id=?", d.Id).
		RunWith(tx).
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
func (d *Donation) GetDonationLink(jg *justgiving.JustGiving) string {
	urls := url.Values{}
	if d.Message != "" {
		urls.Set("message", d.Message)
	}

	urls.Set("currency", d.CurrencyCode)
	urls.Set("amount", AmountToString(d.Amount))
	urls.Set("reference", d.ReferenceCode)

	return jg.GetDonationLink(d.Charity.JustGivingCharityId, urls)
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
	d.LastChecked = pq.NullTime{time.Now(), true}
	err = d.Save(tx)
	return err
}
