package charityhonor

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	TABLE_DRIVES = "drives"
)

var (
	DriveInsertBuilder = QueryBuilder.Insert(TABLE_DRIVES)
)

type Drive struct {
	Id        int `json:"id" db:"id"`
	SourceUrl string
	Amount    int
	Created   time.Time
	Uri       string
}

func (d *Drive) Insert(tx *sqlx.Tx) error {
	if d.Id > 0 {
		return ErrAlreadyInserted
	}

	if err := d.GenerateUri(); err != nil {
		return err
	}

	return DriveInsertBuilder.
		SetMap(d.getSetMap()).
		Suffix(RETURNING_ID).
		RunWith(tx).
		QueryRow().
		Scan(&d.Id)
}

func (d *Drive) getSetMap() M {
	return M{
		"source_url": d.SourceUrl,
		"amount":     d.Amount,
		"uri":        d.Uri,
	}
}

func (d *Drive) GenerateUri() error {
	d.Uri = GenerateUri()
	return nil
}

func (d *Drive) CreateDonation(tx *sqlx.Tx) (*Donation, error) {
	don := &Donation{
		DriveId: d.Id,
	}

	if err := don.GenerateReferenceCode(tx); err != nil {
		return nil, err
	}

	return don, nil
}

func GetDriveByUri(uri string) (*Drive, error) {
	return nil, errors.New("Not implemented")
}
