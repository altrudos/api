package charityhonor

import (
	"errors"
	"os"
)

type M map[string]interface{}

var (
	ErrAlreadyInserted = errors.New("That item has already been inserted into the db")
	ErrNotFound        = errors.New("Couldn't find that")
)

func GetEnv(name, defaultValue string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return defaultValue
}

func GetColumns(colMap map[string]string) []string {
	v := make([]string, 0, len(colMap))
	for _, val := range colMap {
		v = append(v, val)
	}
	return v
}
