package charityhonor

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type Charity struct {
	Id                  int    `db:"id"`
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
	ErrCharityNotFound = errors.New("Charity not found")
)

var (
	CHARITY_COLUMNS = map[string]string{
		"Name":                "name",
		"Description":         "description",
		"JustGivingCharityId": "jg_charity_id",
		"Id":                  "id",
	}
)

func GetCharityById(tx sqlx.Queryer, id int) (*Charity, error) {
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
