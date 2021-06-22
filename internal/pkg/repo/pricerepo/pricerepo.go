package pricerepo

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

const (
	defaultSize    = 64
	defaultTimeout = 1 * time.Second
)

type Model struct {
	ID                 int
	PropertyID         string
	Duration           int
	Amount             int
	Currency           string
	Persons            string
	Weekdays           string
	StayMin            int
	StayMax            int
	XtraPersonPrice    int
	XtraPersonPriceCur string
	PeriodFrom         time.Time
	PeriodTill         time.Time
	Version            uint
}

// ByPersons returns list of person count available for price
func (m *Model) ByPerson() []int {
	list := strings.Split(m.Persons, "|")

	result := make([]int, 0, len(list))
	for _, p := range list {
		personsCount, err := strconv.Atoi(p)
		if err != nil {
			continue
		}
		result = append(result, personsCount)
	}

	return result
}

//var timeLayoutMysql = "2006-01-02 15:04:05"

type DB interface {
	GetDB() *sql.DB
}

type Repo struct {
	db        DB
	tableName string
	log       *zap.Logger
}

func New(db DB, log *zap.Logger, tableName string) *Repo {
	return &Repo{
		db:        db,
		tableName: tableName,
		log:       log.Named("price-repo"),
	}
}

type Filter struct {
	PropertyIDs []string
	PeriodFrom  time.Time
	PeriodTill  time.Time
}

func (r *Repo) List(ctx context.Context, f Filter) ([]Model, error) {
	q := squirrel.
		Select("id",
			"property_id",
			"duration",
			"amount",
			"currency",
			"persons",
			"weekdays",
			"minimum_stay",
			"maximum_stay",
			"extra_person_price",
			"extra_person_price_currency",
			"period_from",
			"period_till",
			"version",
		).
		From(r.tableName)

	if len(f.PropertyIDs) > 0 {
		q = q.Where(squirrel.Eq{"property_id": f.PropertyIDs})
	}

	if !f.PeriodFrom.IsZero() {
		q = q.Where(squirrel.GtOrEq{"period_from": f.PeriodFrom})
	}

	if !f.PeriodTill.IsZero() {
		q = q.Where(squirrel.Lt{"period_till": f.PeriodTill})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	rows, err := r.db.GetDB().QueryContext(dbCtx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			r.log.Error("cannot close db rows")
		}
	}()

	result := make([]Model, 0, defaultSize)
	for rows.Next() {
		var mdl Model
		if err := rows.Scan(
			&mdl.ID,
			&mdl.PropertyID,
			&mdl.Duration,
			&mdl.Amount,
			&mdl.Currency,
			&mdl.Persons,
			&mdl.Weekdays,
			&mdl.StayMin,
			&mdl.StayMax,
			&mdl.XtraPersonPrice,
			&mdl.XtraPersonPriceCur,
			&mdl.PeriodFrom,
			&mdl.PeriodTill,
			&mdl.Version,
		); err != nil {
			return nil, err
		}

		result = append(result, mdl)
	}

	return result, nil
}
