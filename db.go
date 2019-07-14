package charityhonor

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	RETURNING_ID = "RETURNING \"id\""
)

var (
	QueryBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

func MustGetDefaultDb() *sqlx.DB {
	if db != nil {
		return db
	}
	url := GetEnv("PGURL", "postgresql://charityhonor@localhost/charityhonor?sslmode=disable")
	var err error
	db, err = GetPostgresConnection(url)
	if err != nil {
		panic(err)
	}

	return db
}

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
