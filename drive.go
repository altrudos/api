package charityhonor

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	. "github.com/monstercat/pgnull"
)

var (
	TableDrives = "drives"
	ViewDrives  = "drives_view"
)

var DRIVE_COLUMNS map[string]string = map[string]string{
	"Id":        "id",
	"SourceUrl": "source_url",
	"Amount":    "amount",
}

var (
	DriveInsertBuilder = QueryBuilder.Insert(TableDrives)
	DriveSelectBuilder = QueryBuilder.Select(GetColumnsString(DRIVE_COLUMNS)).From(TableDrives)
)

type Drive struct {
	Amount    int
	Created   time.Time
	Id        string `json:"id" setmap:"omitinsert"`
	Source    Source
	SourceUrl string `json:"source_url" db:"source_url"`
	Name      string
	Uri       string

	RedditCommentId NullInt    `db:"reddit_comment_id"`
	RedditUsername  NullString `db:"reddit_username"`
	RedditSubreddit NullString `db:"reddit_subreddit"`
	RedditMarkdown  NullString `db:"reddit_markdown"`

	// From View
	MostRecentDonorAmount int      `db:"most_recent_donor_amount" setmap:"-"`
	MostRecentFinalAmount int      `db:"most_recent_final_amount" setmap:"-"`
	MostRecentTime        NullTime `db:"most_recent_time" setmap:"-"`
	FinalAmountTotal      int      `db:"final_amount_total" setmap:"-"`
	FinalAmountMax        int      `db:"final_amount_max" setmap:"-"`
	DonorAmountTotal      int      `db:"donor_amount_total" setmap:"-"`
	DonorAmountMax        int      `db:"donor_amount_max"`
}

func GetDrives(db sqlx.Queryer, where interface{}) ([]*Drive, error) {
	var xs []*Drive
	c := &Cond{
		Where:    where,
		OrderBys: []string{"created DESC"},
	}
	if err := SelectForStruct(db, &xs, ViewDrives, c); err != nil {
		return nil, err
	}
	return xs, nil
}

func GetDrive(db sqlx.Queryer, where interface{}) (*Drive, error) {
	var x Drive
	if err := GetForStruct(db, &x, ViewDrives, where); err != nil {
		return nil, err
	}
	return &x, nil
}

func GetDriveByUri(uri string) (*Drive, error) {
	return nil, errors.New("Not implemented")
}

func GetDriveById(tx sqlx.Queryer, id string) (*Drive, error) {
	return &Drive{
		Name: "Fake",
	}, nil
}

func GetDriveBySourceUrl(db *sqlx.DB, url string) (*Drive, error) {
	query, args, err := DriveSelectBuilder.
		From(TableDrives).
		Where("source_url=?", url).
		ToSql()
	if err != nil {
		return nil, err
	}
	ds := make([]*Drive, 0)
	err = sqlx.Select(db, &ds, query, args...)
	if err != nil {
		return nil, err
	}
	if len(ds) > 1 {
		return nil, ErrTooManyFound
	}
	if len(ds) == 0 {
		return nil, sql.ErrNoRows
	}

	return ds[0], nil
}

func GetOrCreateDriveBySourceUrl(db *sqlx.DB, url string) (*Drive, error) {
	drive, err := GetDriveBySourceUrl(db, url)
	if err != nil {
		if err == sql.ErrNoRows {
			drive = &Drive{
				SourceUrl: url,
			}
			tx, err := db.Beginx()
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			err = drive.Insert(tx)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			if err := tx.Commit(); err != nil {
				return nil, err
			}

			return drive, nil
		}
		return nil, err
	}

	return drive, nil
}

func (d *Drive) Insert(tx *sqlx.Tx) error {
	if d.Id != "" {
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

func (d *Drive) GenerateDonation() *Donation {
	dono := &Donation{
		DriveId: d.Id,
	}

	return dono
}
