package charityhonor

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	dbUtil "github.com/monstercat/golib/db"
)

const (
	RETURNING_ID = "RETURNING \"id\""
)

var (
	QueryBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

func GetPostgresConnection(url string) (*sqlx.DB, error) {
	connection, err := pq.ParseURL(url)
	if err != nil {
		return nil, err
	}
	db, err := sqlx.Open("postgres", connection)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetColumns(colMap map[string]string) []string {
	v := make([]string, 0, len(colMap))
	for _, val := range colMap {
		v = append(v, val)
	}
	return v
}

//Turns a map of strings into a "id, name, title" type string
//The map values are used, the keys are ignored
func GetColumnsString(colMap map[string]string) string {
	cols := GetColumns(colMap)

	return strings.Join(cols, ", ")
}

func GetColumnsTabled(colMap map[string]string, table string) []string {
	v := make([]string, 0, len(colMap))
	for _, val := range colMap {
		v = append(v, table+"."+val)
	}
	return v
}

type Cond struct {
	Where    interface{}
	OrderBys []string
	Limit    int
	Offset   int
}

func (c *Cond) DefaultLimit(limit int) {
	if c.Limit == 0 {
		c.Limit = limit
	}
}
func (c *Cond) DefaultOffset(offset int) {
	if c.Offset == 0 {
		c.Offset = offset
	}
}

func (c *Cond) DefaultOrderBys(orders ...string) {
	if len(c.OrderBys) == 0 {
		c.OrderBys = orders
	}
}

func (c *Cond) ApplyWhere(qry *squirrel.SelectBuilder) {
	if c.Where != nil {
		*qry = qry.Where(c.Where)
	}
}

func (c *Cond) ApplyWithoutLimits(qry *squirrel.SelectBuilder) {
	c.ApplyWhere(qry)
	if len(c.OrderBys) > 0 {
		*qry = qry.OrderBy(c.OrderBys...)
	}
}
func (c *Cond) ApplyLimits(qry *squirrel.SelectBuilder) {
	if c.Limit > 0 {
		*qry = qry.Limit(uint64(c.Limit))
	}
	if c.Offset > 0 {
		*qry = qry.Offset(uint64(c.Offset))
	}
}

func (c *Cond) Apply(qry *squirrel.SelectBuilder) {
	c.ApplyWithoutLimits(qry)
	c.ApplyLimits(qry)
}

func SelectForStruct(db sqlx.Queryer, slice interface{}, table string, cond *Cond) error {
	cols := dbUtil.GetColumnsList(slice, "")
	qry := QueryBuilder.Select(cols...).From(table)
	if cond != nil {
		cond.Apply(&qry)
	}
	return dbUtil.Select(db, slice, qry)
}

func GetForStruct(db sqlx.Queryer, val interface{}, table string, where interface{}) error {
	cols := dbUtil.GetColumnsList(val, "")
	qry := QueryBuilder.Select(cols...).From(table)
	if where != nil {
		qry = qry.Where(where)
	}
	return dbUtil.Get(db, val, qry)

}

func GetCount(db sqlx.Queryer, qry squirrel.SelectBuilder) (int, error) {
	var total int
	if err:= dbUtil.Get(db, &total, qry); err != nil {
		return 0, err
	}
	return total, nil
}