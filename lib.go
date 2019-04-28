package charityhonor

import (
	"errors"
	"os"
	"strconv"
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

func AmountToString(amount int) string {
	str := strconv.Itoa(amount)

	if len(str) == 1 {
		str = "0" + str
	}

	first := str[:len(str)-2]
	last := str[len(str)-2:]

	if first == "" {
		first = "0"
	}

	if len(last) == 1 {
		last = last + "0"
	} else if len(last) == 0 {
		last = "00"
	}

	return first + "." + last
}
