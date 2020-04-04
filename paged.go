package charityhonor

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	dbUtil "github.com/monstercat/golib/db"
)

type QueryGenerator func(cols... string) squirrel.SelectBuilder

type Paged struct {
	Data   interface{}
	Limit  int
	Offset int
	Total  int
}

func GetWithTotal(
	db sqlx.Queryer,
	generator QueryGenerator,
	slice interface{},
	cond *Cond,
) (int, error) {

	total, err := GetCount(db, generator("COUNT(*)"))
	if err != nil {
		return 0, err
	}

	cols := dbUtil.GetColumnsList(slice, "")
	qry := generator(cols...)
	cond.ApplyLimits(&qry)
	if err := dbUtil.Select(db, slice, qry); err != nil {
		return 0, err
	}

	return total, nil
}

var AggregateFunctionList = []string{
	"array_agg",
	"avg",
	"bit_and",
	"bit_or",
	"bool_and",
	"bool_or",
	"count",
	"every",
	"json_agg",
	"jsonb_agg",
	"json_object_agg",
	"jsonb_object_agg",
	"max",
	"min",
	"string_agg",
	"sum",
	"xmlagg",

	// Aggregate for stats
	"corr",
	"covar_pop",
	"covar_samp",
	"regr_avgx",
	"regr_avgy",
	"regr_count",
	"regr_intercept",
	"regr_r2",
	"regr_slope",
	"regr_sxx",
	"regr_sxy",
	"regr_syy",
	"stddev",
	"stddev_pop",
	"stddev_samp",
	"variance",
	"var_pop",
	"var_samp",

	// ordered-set agg
	"mode",
	"percentile_cont",
	"percentile_disc",

	// hypothetical-set
	"rank", // also handles dense_rank & percent_rank
	"cume_dist",

	// grouping
	"grouping",
}

func hasAggregate(cols []string) bool {
	for _, c := range cols {
		for _, f := range AggregateFunctionList {
			if strings.Index(strings.ToLower(c), f + "(") > -1 {
				return true
			}
		}
	}
	return false
}

func DefaultGenerator(
	table string,
	cond *Cond,
) QueryGenerator {
	return func(cols...string) squirrel.SelectBuilder {
		qry := QueryBuilder.Select(cols...).From(table)
		if hasAggregate(cols) {
			cond.ApplyWhere(&qry)
		} else {
			cond.ApplyWithoutLimits(&qry)
		}
		return qry
	}
}