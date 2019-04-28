package charityhonor

import (
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
