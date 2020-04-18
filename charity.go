package charityhonor

import (
	"errors"
	"sort"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/monstercat/golib/db"
	. "github.com/monstercat/pgnull"

	"github.com/charityhonor/ch-api/pkg/justgiving"

	"github.com/jmoiron/sqlx"
)

type Charity struct {
	CountryCode         string `db:"country_code"`
	Description         string `db:"description"`
	FeatureScore        int    `db:"feature_score"`
	Id                  string `setmap:"omitinsert"`
	JustGivingCharityId int    `db:"jg_charity_id"`
	Name                string `db:"name"`
	LogoUrl             string `db:"logo_url"`
	Subtext             string `db:"subtext"`
	WebsiteUrl          string `db:"website_url"`

	// From View
	MostRecentDonorAmount int      `db:"most_recent_donor_amount" setmap:"-"`
	MostRecentUSDAmount   int      `db:"most_recent_usd_amount" setmap:"-"`
	MostRecentTime        NullTime `db:"most_recent_time" setmap:"-"`
	USDAmountTotal        int      `db:"usd_amount_total" setmap:"-"`
	DonorAmountTotal      int      `db:"donor_amount_total" setmap:"-"`

	// Filled in afterwards
	Top10Donations []*Donation `db:"-" setmap:"-"`
}

var (
	TableCharities        = "charities"
	ViewCharities         = "charities_view"
	ViewDonations         = "donations_view"
	ViewFeaturedCharities = "featured_charities_view"
)

var (
	CharityInsertBuilder = QueryBuilder.Insert(TableCharities)
	CharitySelectBuilder = QueryBuilder.Select(GetColumns(CharityColumns)...).From(TableCharities)
)

var (
	ErrCharityNotFound      = errors.New("Charity not found")
	ErrDuplicateJGCharityId = errors.New("a charity with that JustGiving charity ID already exists")
)

var (
	CharityColumns = map[string]string{
		"Name":                "name",
		"Description":         "description",
		"FeatureScore":        "feature_score",
		"Subtext":             "subtext",
		"LogoUrl":             "logo_url",
		"CountryCode":         "country_code",
		"JustGivingCharityId": "jg_charity_id",
		"Id":                  "id",
	}
)

func ConvertCharityError(err error) error {
	if err == nil {
		return nil
	}

	if ErrIsPqConstraint(err, "charities_jg_charity_id_unique") {
		return ErrDuplicateJGCharityId

	}

	return err
}

func (c *Charity) Insert(ext sqlx.Ext) error {
	smap := dbUtil.SetMap(c, true)
	err := CharityInsertBuilder.
		SetMap(smap).
		Suffix(PqSuffixId).
		RunWith(ext).
		QueryRow().
		Scan(&c.Id)

	return ConvertCharityError(err)
}

func GetCharityTop10Donations(db sqlx.Queryer, cId string) ([]*Donation, error) {
	var xs []*Donation
	cond := &Cond{
		Where:    squirrel.Eq{"charity_id": cId},
		OrderBys: []string{"-usd_amount"},
		Limit:    10,
	}
	if err := SelectForStruct(db, &xs, ViewDonations, cond); err != nil {
		return nil, err
	}
	return xs, nil
}

func GetCharityById(tx sqlx.Queryer, id string) (*Charity, error) {
	query, args, err := CharitySelectBuilder.
		Where("id=?", id).
		ToSql()
	if err != nil {
		return nil, err
	}

	var d Charity
	err = sqlx.Get(tx, &d, query, args...)
	return &d, err
}

func GetCharityByName(tx sqlx.Queryer, name string) (*Charity, error) {
	query, args, err := CharitySelectBuilder.
		Where("LOWER(name)=?", strings.ToLower(name)).
		ToSql()
	if err != nil {
		return nil, err
	}

	var d Charity
	err = sqlx.Get(tx, &d, query, args...)
	return &d, err
}

func GetCharities(db sqlx.Queryer, cond *Cond) (xs []*Charity, err error) {
	err = SelectForStruct(db, &xs, ViewCharities, cond)
	return
}

func GetCharity(db sqlx.Queryer, cond *Cond) (*Charity, error) {
	var x Charity
	err := GetForStruct(db, &x, ViewCharities, cond)
	return &x, err
}

func GetCharitiesByJGId(db sqlx.Queryer, ids pq.Int64Array) ([]*Charity, error) {
	cond := &Cond{
		Where:    squirrel.Expr("jg_charity_id=ANY(?)", ids),
		OrderBys: []string{"feature_score DESC", "name ASC"},
	}
	return GetCharities(db, cond)
}

func GetCharityByJGId(db sqlx.Queryer, id int) (*Charity, error) {
	cond := &Cond{
		Where: squirrel.Expr("jg_charity_id=", id),
	}
	return GetCharity(db, cond)
}

func ConvertJGCharity(jgc *justgiving.Charity) *Charity {
	return &Charity{
		Name:                jgc.Name,
		Description:         jgc.Description,
		JustGivingCharityId: jgc.Id,
	}
}

func ConvertJGCharities(jgcs []*justgiving.Charity) []*Charity {
	charities := make([]*Charity, len(jgcs))
	for i, v := range jgcs {
		charities[i] = ConvertJGCharity(v)
	}
	return charities
}

func SortCharities(charities []*Charity) {
	sort.Slice(charities, func(i, j int) bool {
		a := charities[i]
		b := charities[j]
		if a.FeatureScore > b.FeatureScore {
			return true
		}
		return a.Name < b.Name
	})
}
