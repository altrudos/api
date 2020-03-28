package charityhonor

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func GetTestDb() *sqlx.DB {
	if db != nil {
		return db
	}
	url := GetEnv("TESTPGURL", "postgresql://charityhonor@localhost/ch?sslmode=disable")
	var err error
	db, err = GetPostgresConnection(url)
	if err != nil {
		fmt.Println("NO CONNECTION")
		panic(err)
	}

	return db
}
