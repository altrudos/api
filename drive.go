package charityhonor

import (
	"database/sql"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"

	"github.com/monstercat/golib/db"

	"github.com/jmoiron/sqlx"
	. "github.com/monstercat/pgnull"
)

var (
	TableDrives = "drives"
	ViewDrives  = "drives_view"
)

var (
	DriveInsertBuilder = QueryBuilder.Insert(TableDrives)
	DriveSelectBuilder = QueryBuilder.Select("*").From(ViewDrives)
)

type Drive struct {
	Amount     int
	Created    time.Time
	Id         string     `setmap:"omitinsert"`
	Source     Source     `db:"-"`
	SourceUrl  string     `db:"source_url"`
	SourceKey  string     `db:"source_key"`
	SourceType SourceType `db:"source_type"`
	SourceMeta FlatMap    `db:"source_meta"`
	Uri        string

	// From View
	MostRecentDonorAmount int      `db:"most_recent_donor_amount" setmap:"-"`
	MostRecentFinalAmount int      `db:"most_recent_final_amount" setmap:"-"`
	MostRecentTime        NullTime `db:"most_recent_time" setmap:"-"`
	FinalAmountTotal      int      `db:"final_amount_total" setmap:"-"`
	FinalAmountMax        int      `db:"final_amount_max" setmap:"-"`
	DonorAmountTotal      int      `db:"donor_amount_total" setmap:"-"`
	DonorAmountMax        int      `db:"donor_amount_max" setmap:"-"`

	// Filled in afterwards
	Top10Donations    []*Donation `db:"-" setmap:"-"`
	Recent10Donations []*Donation `db:"-" setmap:"-"`
}

// For queries that include the TopAmount and NumDonations calculations
type DriveTallied struct {
	Drive
	// From sum queries
	TopAmount    int `db:"top_amount" setmap:"-"`
	NumDonations int `db:"num_donations" setmap:"-"`
}

func GetDrive(db sqlx.Queryer, where interface{}) (*Drive, error) {
	var x Drive
	if err := GetForStruct(db, &x, ViewDrives, where); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &x, nil
}

/*SELECT pp.id
     , count(kind = 'dog' OR NULL) AS alive_dogs_count
     , count(kind = 'cat' OR NULL) AS alive_cats_count
FROM   people pp
LEFT   JOIN pets pt ON pt.person_id = pp.id
                   AND pt.alive
WHERE  <some condition to retrieve a small subset>
GROUP  BY 1;*/
func GetTopDrives(db sqlx.Queryer) ([]*DriveTallied, error) {
	qry := QueryBuilder.Select("dr.*, sq.top_amount").
		From(`(SELECT SUM(final_amount) as top_amount, COUNT(dono.id) as num_donations, drive_id
		FROM ` + TableDonations + ` dono
		WHERE dono.created >= NOW() - INTERVAL '7 DAYS' 
		AND dono.status = 'Accepted'
		GROUP BY drive_id) sq`).
		Join(ViewDrives + " dr ON dr.id = sq.drive_id").
		OrderBy("top_amount DESC")

	var drives []*DriveTallied
	if err := dbUtil.Select(db, &drives, qry); err != nil {
		return nil, err
	}
	return drives, nil
}

func GetDriveByField(q sqlx.Queryer, field, value string) (*Drive, error) {
	return GetDrive(q, squirrel.Eq{field: value})
}

func GetDriveByUri(q sqlx.Queryer, uri string) (*Drive, error) {
	return GetDriveByField(q, "LOWER(uri)", strings.ToLower(uri))
}

func GetDriveById(q sqlx.Queryer, id string) (*Drive, error) {
	return GetDriveByField(q, "id", id)
}

func GetDriveTopDonations(db sqlx.Queryer, cId string, num int) ([]*Donation, error) {
	var xs []*Donation
	cond := &Cond{
		Where: squirrel.Eq{
			"drive_id": cId,
			"status":   DonationAccepted,
		},
		OrderBys: []string{"-final_amount"},
		Limit:    num,
	}
	if err := SelectForStruct(db, &xs, TableDonations, cond); err != nil {
		return nil, err
	}
	return xs, nil
}

func GetDriveRecentDonations(db sqlx.Queryer, cId string, num int) ([]*Donation, error) {
	var xs []*Donation
	cond := &Cond{
		Where: squirrel.Eq{
			"drive_id": cId,
			"status":   DonationAccepted,
		},
		OrderBys: []string{"created DESC"},
		Limit:    num,
	}
	if err := SelectForStruct(db, &xs, TableDonations, cond); err != nil {
		return nil, err
	}
	return xs, nil
}

func GetOrCreateDriveBySourceUrl(ext sqlx.Ext, url string) (*Drive, error) {
	drive, err := GetDriveBySourceUrl(ext, url)
	if err == nil {
		return drive, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}
	drive, err = CreatedDriveBySourceUrl(ext, url)
	if err != nil {
		return nil, err
	}
	err = drive.Create(ext)
	if err != nil {
		return nil, err
	}

	return drive, nil
}

func GetDriveBySourceUrl(q sqlx.Queryer, url string) (*Drive, error) {
	source, err := ParseSourceURL(url)
	if err != nil {
		return nil, err
	}
	return GetDriveBySource(q, source)
}

func GetDriveBySource(q sqlx.Queryer, source Source) (*Drive, error) {
	eq := squirrel.Eq{
		"source_type": source.GetType(),
		"source_key":  source.GetKey(),
	}
	return GetDrive(q, eq)
}
func CreatedDriveBySourceUrl(ext sqlx.Ext, url string) (*Drive, error) {
	source, err := ParseSourceURL(url)
	if err != nil {
		return nil, err
	}
	meta, err := source.GetMeta()
	if err != nil {
		return nil, err
	}
	drive := &Drive{
		Source:     source,
		SourceUrl:  url,
		SourceKey:  source.GetKey(),
		SourceType: source.GetType(),
		SourceMeta: meta,
	}

	if err := drive.Create(ext); err != nil {
		return nil, err
	}

	return drive, nil
}

func (d *Drive) Create(ext sqlx.Ext) error {
	if d.Id != "" {
		return ErrAlreadyInserted
	}
	if d.Uri == "" {
		if err := d.GenerateUri(); err != nil {
			return err
		}
	}
	d.Created = time.Now()
	return d.Insert(ext)
}

func (d *Drive) Insert(ext sqlx.Ext) error {
	setMap := dbUtil.SetMap(d, true)
	query := DriveInsertBuilder.
		SetMap(setMap).
		Suffix(RETURNING_ID)

	return query.
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

func (d *Drive) GetDonationQueryBuilder() squirrel.SelectBuilder {
	return QueryBuilder.Select(dbUtil.GetColumnsList(&Donation{}, "")...).
		From(TableDonations).
		Where("drive_id=?", d.Id)
}

func (d *Drive) ApprovedDonations(q sqlx.Queryer, limit int) *squirrel.SelectBuilder {
	query := d.GetDonationQueryBuilder()

	if limit > 0 {
		query = query.Limit(uint64(limit))
	} else {
		query = query.Limit(5)
	}

	ApplyApproved(&query)
	return &query
}

func (d *Drive) GetTopDonations(q sqlx.Queryer, limit int) ([]*Donation, error) {
	query := d.ApprovedDonations(q, limit)
	*query = query.OrderBy("final_amount DESC")

	return QueryDonations(q, query)
}

func (d *Drive) GetRecentDonations(q sqlx.Queryer, limit int) ([]*Donation, error) {
	query := d.ApprovedDonations(q, limit)
	*query = query.OrderBy("created DESC")

	return QueryDonations(q, query)
}

func GetDrives(db sqlx.Queryer, cond *Cond) (xs []*Drive, err error) {
	err = SelectForStruct(db, &xs, ViewDrives, cond)
	return
}
