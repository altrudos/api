package charityhonor

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	pgUtils "github.com/monstercat/golib/db/postgres"
)

const TableSearchCache = "search_cache"

type SearchCacheItem struct {
	Term    string        `db:"term"`
	Expires time.Time     `db:"expires"`
	Ids     pq.Int64Array `db:"ids"`
}

func GetSearchItem(db sqlx.Queryer, term string) (*SearchCacheItem, error) {
	var item SearchCacheItem
	if err := GetForStruct(db, &item, TableSearchCache, squirrel.Expr("term ILIKE ?", term)); err != nil {
		return nil, err
	}
	return &item, nil
}

func CreateSearchCache(db sqlx.Ext, term string, ids pq.Int64Array) error {
	item := &SearchCacheItem{
		Term:    term,
		Ids:     ids,
		Expires: time.Now().Add(7 * 24 * time.Hour),
	}
	return item.Insert(db)
}

func (i *SearchCacheItem) Expired() bool {
	return i.Expires.Before(time.Now())
}

func (i *SearchCacheItem) Insert(db sqlx.Ext) error {
	return pgUtils.InsertSetMapNoId(db, TableSearchCache, i)
}

func (i *SearchCacheItem) Delete(db sqlx.Ext) error {
	return pgUtils.DeleteWhere(db, TableSearchCache, squirrel.Eq{"term": i.Term})
}