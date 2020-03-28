package charityhonor

import (
	"database/sql"
	"strings"
	"time"

	"github.com/monstercat/golib/db"

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
	Id        string  `json:"id" setmap:"omitinsert"`
	Source    *Source `db:"-"`
	SourceUrl string  `json:"source_url" db:"source_url"`
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

func GetDriveByField(q sqlx.Queryer, field, value string) (*Drive, error) {
	query, args, err := DriveSelectBuilder.
		From(TableDrives).
		Where(field+"=?", value).
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

func GetDriveByUri(q sqlx.Queryer, uri string) (*Drive, error) {
	return GetDriveByField(q, "LOWER(uri)", strings.ToLower(uri))
}

func GetDriveById(q sqlx.Queryer, id string) (*Drive, error) {
	return GetDriveByField(q, "id", id)

}

func GetDriveBySourceUrl(q sqlx.Queryer, url string) (*Drive, error) {
	return GetDriveByField(q, "source_url", url)
}

func GetOrCreateDriveBySourceUrl(ext sqlx.Ext, url string) (*Drive, error) {
	drive, err := GetDriveBySourceUrl(db, url)
	if err == nil {
		return drive, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}
	drive = &Drive{
		SourceUrl: url,
	}
	err = drive.Create(ext)
	if err != nil {
		return nil, err
	}

	return drive, nil
}

func (d *Drive) Create(ext sqlx.Ext) error {
	if d.Id != "" {
		return ErrAlreadyInserted
	}
	if err := d.GenerateUri(); err != nil {
		return err
	}
	return d.Insert(ext)
}

func (d *Drive) Insert(ext sqlx.Ext) error {
	return DriveInsertBuilder.
		SetMap(dbUtil.SetMap(d, true)).
		Suffix(RETURNING_ID).
		RunWith(ext).
		QueryRow().
		Scan(&d.Id)
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
