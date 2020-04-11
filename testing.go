package charityhonor

import (
	"github.com/jmoiron/sqlx"
)

func MustGetTestConfig() *Config {
	c, err := ParseConfig(GetEnv("TESTCONFIG", "config_test.toml"))
	if err != nil {
		panic(err)
	}
	return c
}

func GetTestServices() *Services {
	services, err := MustGetTestConfig().Connect()
	if err != nil {
		panic(err)
	}
	return services
}

func GetTestDb() *sqlx.DB {
	services := GetTestServices()
	return services.DB
}
