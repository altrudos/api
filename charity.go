package charityhonor

import (
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/monstercat/golib/db"

	"github.com/charityhonor/ch-api/pkg/justgiving"

	"github.com/jmoiron/sqlx"
)

type Charity struct {
	Id                  string `db:"id"`
	Name                string `db:"name"`
	Description         string `db:"description"`
	JustGivingCharityId int    `db:"jg_charity_id"`
}

var (
	TABLE_CHARITIES = "charities"
)

var (
	CharityInsertBuilder = QueryBuilder.Insert(TABLE_CHARITIES)
	CharitySelectBuilder = QueryBuilder.Select(GetColumns(CHARITY_COLUMNS)...).From(TABLE_CHARITIES)
)

var (
	ErrCharityNotFound      = errors.New("Charity not found")
	ErrDuplicateJGCharityId = errors.New("a charity with that JustGiving charity ID already exists")
)

var (
	CHARITY_COLUMNS = map[string]string{
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
	if err, ok := err.(*pq.Error); ok {
		// Here err is of type *pq.Error, you may inspect all its fields, e.g.:
		name := err.Code.Name()
		fmt.Println("name", name)
		return ErrDuplicateJGCharityId
	}
}

func (c *Charity) Insert(ext sqlx.Ext) error {
	_, err := CharityInsertBuilder.
		Values(dbUtil.SetMap(c, true)).
		Suffix(PqSuffixId).
		RunWith(ext).
		Exec()

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
