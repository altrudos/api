package charityhonor

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/jmoiron/sqlx"

	"github.com/charityhonor/ch-api/pkg/justgiving"
)

var (
	ErrInvalidJGMode = errors.New("invalid justgiving mode")
)

type Config struct {
	Postgres   string
	JustGiving JGConfig
}

type JGConfig struct {
	AppId string
	Mode  justgiving.Mode
}

type Services struct {
	DB *sqlx.DB
	JG *justgiving.JustGiving
}

func MustGetConfigServices(confFile string) *Services {
	s, err := GetConfigServices(confFile)
	if err != nil {
		panic(err)
	}
	return s
}

func GetConfigServices(confFile string) (*Services, error) {
	c, err := ParseConfig(confFile)
	if err != nil {
		return nil, err
	}
	return c.Connect()
}

func ParseConfig(confFile string) (*Config, error) {
	f, err := os.Open(confFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var conf Config
	if _, err := toml.DecodeReader(f, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func (c *Config) MustConnect() *Services {
	s, err := c.Connect()
	if err != nil {
		panic(err)
	}
	return s
}

func (c *Config) Connect() (*Services, error) {
	db, err := c.ConnectPostgres()
	if err != nil {
		return nil, err
	}
	jg, err := c.ConnectJG()
	if err != nil {
		return nil, err
	}
	return &Services{
		DB: db,
		JG: jg,
	}, nil
}

func (c *Config) MustConnectPostgres() *sqlx.DB {
	db, err := c.ConnectPostgres()
	if err != nil {
		panic(err)
	}
	return db
}

func (c *Config) ConnectPostgres() (*sqlx.DB, error) {
	return GetPostgresConnection(c.Postgres)
}

func (c *Config) ConnectJG() (*justgiving.JustGiving, error) {
	return c.JustGiving.Connect()
}

func (c *JGConfig) Connect() (*justgiving.JustGiving, error) {
	// Make sure mode is valid.
	validModes := []justgiving.Mode{justgiving.ModeProduction, justgiving.ModeStaging}
	var isValid bool
	for _, v := range validModes {
		if c.Mode == v {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, ErrInvalidJGMode
	}

	return &justgiving.JustGiving{
		AppId: c.AppId,
		Mode:  justgiving.Mode(c.Mode),
	}, nil
}
