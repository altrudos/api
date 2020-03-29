package charityhonor

import (
	"errors"
	"strings"

	"github.com/monstercat/golib/db"
	. "github.com/monstercat/pgnull"

	"github.com/charityhonor/ch-api/pkg/justgiving"

	"github.com/jmoiron/sqlx"
)

type Charity struct {
	Description         string
	Id                  string `setmap:"omitinsert"`
	JustGivingCharityId int    `db:"jg_charity_id"`
	Name                string

	// From View
	MostRecentDonorAmount int      `db:"most_recent_donor_amount" setmap:"-"`
	MostRecentFinalAmount int      `db:"most_recent_final_amount" setmap:"-"`
	MostRecentTime        NullTime `db:"most_recent_time" setmap:"-"`
	FinalAmountTotal      int      `db:"final_amount_total" setmap:"-"`
	FinalAmountMax        int      `db:"final_amount_max" setmap:"-"`
	DonorAmountTotal      int      `db:"donor_amount_total" setmap:"-"`
	DonorAmountMax        int      `db:"donor_amount_max" setmap:"-"`
}

var (
	TableCharities = "charities"
	ViewCharities  = "charities_view"
)

var (
	CharityInsertBuilder = QueryBuilder.Insert(TableCharities)
	CharitySelectBuilder = QueryBuilder.Select(GetColumns(CharityColumns)...).From(TableCharities)
)

var (
	ErrCharityNotFound      = errors.New("Charity not found")
	ErrDuplicateJGCharityId = errors.New("a charity with that JustGiving charity ID already exists")
)

var (
	CharityColumns = map[string]string{
		"Name":                "name",
		"Description":         "description",
		"JustGivingCharityId": "jg_charity_id",
		"Id":                  "id",
	}
)

func ConvertCharityError(err error) error {
	if err == nil {
		return nil
	}

	if ErrIsPqConstraint(err, "charities_jg_charity_id_unique") {
		return ErrDuplicateJGCharityId

	}

	return err
}

func (c *Charity) Insert(ext sqlx.Ext) error {
	err := CharityInsertBuilder.
		SetMap(dbUtil.SetMap(c, true)).
		Suffix(PqSuffixId).
		RunWith(ext).
		QueryRow().
		Scan(&c.Id)

	return ConvertCharityError(err)
}

func GetCharityById(tx sqlx.Queryer, id string) (*Charity, error) {
	query, args, err := CharitySelectBuilder.
		Where("id=?", id).
		ToSql()
	if err != nil {
		return nil, err
	}

	var d Charity
	err = sqlx.Get(tx, &d, query, args...)
	return &d, err
}

func GetCharityByName(tx sqlx.Queryer, name string) (*Charity, error) {
	query, args, err := CharitySelectBuilder.
		Where("LOWER(name)=?", strings.ToLower(name)).
		ToSql()
	if err != nil {
		return nil, err
	}

	var d Charity
	err = sqlx.Get(tx, &d, query, args...)
	return &d, err
}

func ConvertJGCharity(jgc *justgiving.Charity) *Charity {
	return &Charity{
		Name:                jgc.Name,
		Description:         jgc.Description,
		JustGivingCharityId: jgc.Id,
	}
}

func ConvertJGCharities(jgcs []*justgiving.Charity) []*Charity {
	charities := make([]*Charity, len(jgcs))
	for i, v := range jgcs {
		charities[i] = ConvertJGCharity(v)
	}
	return charities
}
