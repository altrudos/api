package charityhonor

import (
	"errors"
	"fmt"
	"os"

	"github.com/lib/pq"
)

type M map[string]interface{}

var (
	ErrAlreadyInserted = errors.New("That item has already been inserted into the db")
	ErrNotFound        = errors.New("Couldn't find that")
	ErrTooManyFound    = errors.New("Found too many of that")
)

const PqSuffixId = "RETURNING \"id\""

type SortDirection string

var (
	SortDesc = SortDirection("DESC")
	SortAsc  = SortDirection("ASC")
)

type QueryOperator interface {
	GetSort() string
	GetLimit() int
}

type BaseOperator struct {
	SortField string
	SortDir   SortDirection
	Limit     int
}

func (qo *BaseOperator) GetLimit() int {
	if qo.Limit == 0 {
		return 50
	}
	return qo.Limit
}

func (qo *BaseOperator) GetSort() string {
	if qo.SortField == "" {
		return ""
	}

	direction := qo.SortDir
	if direction == "" {
		direction = SortAsc
	}

	return fmt.Sprintf("%s %s", qo.SortField, direction)
}

func StatusesPQStringArray(things []DonationStatus) pq.StringArray {
	strs := make([]string, len(things))
	for i, thing := range things {
		strs[i] = fmt.Sprintf("%s", thing)
	}
	return pq.StringArray(strs)
}

func GetEnv(name, defaultValue string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return defaultValue
}

func AmountToString(amount int) string {
	str := fmt.Sprintf("%.2f", float64(amount)/100.0)
	return str
}

func ErrIsPqConstraint(err error, constraint string) bool {
	if err, ok := err.(*pq.Error); ok {
		if err.Constraint == constraint {
			return true
		}
	}

	return false
}
