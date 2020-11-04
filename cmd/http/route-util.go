package main

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	. "github.com/altrudos/api"
)

type Getter func(db sqlx.Queryer, id string) (interface{}, error)

func getById(
	paramKey string,
	itemName string,
	fn Getter,
) HandlerFunc {
	return func(c *RouteContext) {
		id := c.Params[paramKey]
		if c.HandledMissingParam(id) {
			return
		}

		item, err := fn(c.DB, id)
		if c.HandledError(err) {
			return
		}
		c.JSON(http.StatusOK, M{
			itemName: item,
		})
	}
}

func defaultGetAll(
	c *RouteContext,
	itemName string,
	tableName string,
	slice interface{},
	cond *Cond,
) {
	generator := DefaultGenerator(tableName, cond)
	getAll(c, itemName, generator, slice, cond)
}

func getAll(
	c * RouteContext,
	itemName string,
	generator QueryGenerator,
	slice interface{},
	cond *Cond,
) {
	total, err := GetWithTotal(c.DB, generator, slice, cond)
	if c.HandledError(err) {
		return
	}
	c.JSON(http.StatusOK, M{
		itemName: Paged{
			Data:   slice,
			Total:  total,
			Offset: cond.Offset,
			Limit:  cond.Limit,
		},
	})

}
