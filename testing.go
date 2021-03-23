package altrudos

import (
	"github.com/jmoiron/sqlx"
)

func MustGetTestConfig() *Config {
	filepath := GetEnv("ALTRUDOS_TESTCONFIG", "config_test.toml")
	c, err := ParseConfig(filepath)
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
