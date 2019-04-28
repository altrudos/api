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
	url := GetEnv("TESTPGURL", "postgresql://charityhonor@localhost/charityhonor?sslmode=disable")
	fmt.Println("url", url)
	var err error
	db, err = GetPostgresConnection(url)
	if err != nil {
		panic(err)
	}

	return db
}
