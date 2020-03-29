package charityhonor

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	dbUtil "github.com/monstercat/golib/db"
)

type QueryGenerator func(cols... string) squirrel.SelectBuilder

type Paged struct {
	Data   interface{}
	Limit  int
	Offset int
	Total  int
}

func GetWithTotal(
	db sqlx.Queryer,
	generator QueryGenerator,
	slice interface{},
	cond *Cond,
) (int, error) {

	total, err := GetCount(db, generator("COUNT(*)"))
	if err != nil {
		return 0, err
	}

	//TODO: we need to update dbutils to handle slices properly.
	// Then we can call GetColumnsList on a slice to get the columns.
	// And not have to use *
	qry := generator("*")
	cond.ApplyLimits(&qry)
	if err := dbUtil.Select(db, slice, qry); err != nil {
		return 0, err
	}

	return total, nil
}

func DefaultGenerator(
	table string,
	cond *Cond,
) QueryGenerator {
	return func(cols...string) squirrel.SelectBuilder {
		qry := QueryBuilder.Select(cols...).From(table)
		cond.ApplyWithoutLimits(&qry)
		return qry
	}
}