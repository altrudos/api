package charityhonor

import (
	"github.com/jmoiron/sqlx"
)

func GetTestServices() *Services {
	return MustGetConfigServices(GetEnv("TESTCONFIG", "config_test.toml"))
}

func GetTestDb() *sqlx.DB {
	services := GetTestServices()
	return services.DB
}
