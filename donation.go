package charityhonor

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
)

var (
	DONATION_COLUMNS = map[string]string{
		"LastChecked":   "last_checked",
		"Status":        "status",
		"ReferenceCode": "reference_code",
		"DriveId":       "drive_id",
		"Created":       "created",
		"DonorName":     "donor_name",
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
	ReferenceCode string `db:"reference_code"`
	DriveId       int    `db:"drive_id"`
	CharityId     int    `db:"charity_id"`
	Amount        int
	DonorName     string `db:"donor_name"`
	Message       string
	CurrencyCode  string `db:"currency_code"`
}

func (d *Donation) GenerateReferenceCode(tx *sqlx.Tx) error {
	exists := false
	for d.ReferenceCode == "" || exists == true {
		str := time.Now().String()
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

func (d *Donation) Insert(tx *sqlx.Tx) error {
	if d.ReferenceCode == "" {
		return ErrMissingReferenceCode
	}
	return DonationInsertBuilder.
		SetMap(d.getSetMap()).
		Suffix(RETURNING_ID).
		RunWith(tx).
		QueryRow().
		Scan(&d.Id)
}

func GetDonationByReferenceCode(tx sqlx.Queryer, code string) (*Donation, error) {
	query, args, err := QueryBuilder.
		Select(GetColumns(DONATION_COLUMNS)...).
		From(TABLE_DONATIONS).Where("reference_code=?", code).
		ToSql()
	if err != nil {
		return nil, err
	}

	var d Donation
	err = sqlx.Get(tx, &d, query, args...)
	fmt.Println("err", err)
	return &d, err
}
