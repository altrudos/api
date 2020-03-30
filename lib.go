package charityhonor

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/lib/pq"
)

type M map[string]interface{}

type FlatMap map[string]interface{}

func (p FlatMap) Equal(b FlatMap) bool {
	for k, v := range p {
		e, ok := b[k]
		if !ok {
			return false
		}

		if !reflect.DeepEqual(e, v) {
			return false
		}
	}
	return true
}

func (p FlatMap) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

func (p *FlatMap) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	// Nothing to do if nil
	if i == nil {
		return nil
	}

	*p, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("Type assertion .(map[string]interface{}) failed.")
	}

	return nil
}

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
