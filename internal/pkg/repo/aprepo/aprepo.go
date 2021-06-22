package aprepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

const (
	defaultSize    = 64
	defaultTimeout = 1 * time.Second
)

type Model struct {
	PropertyID       string
	Date             time.Time
	Quantity         int
	ArrivalAllowed   bool
	DepartureAllowed bool
	ArrivalStayMin   int
	ArrivalStayMax   int

	Duration           int
	Amount             int
	Currency           string
	Persons            string
	Weekdays           string
	PriceStayMin       int
	PriceStayMax       int
	XtraPersonPrice    int
	XtraPersonPriceCur string
}

const timeLayoutMysql = "2006-01-02"

//const timeLayoutMysql = "2006-01-02 15:04:05"

type DB interface {
	GetDB() *sql.DB
}

type Repo struct {
	db  DB
	log *zap.Logger
}

func New(db DB, log *zap.Logger) *Repo {
	return &Repo{
		db:  db,
		log: log.Named("availabilities-repo"),
	}
}

type Filter struct {
	PropertyIDs []string
	DateFrom    time.Time
	DateTill    time.Time
}

func (r *Repo) List(ctx context.Context, f Filter) ([]Model, error) {
	if r.db == nil {
		return nil, errors.New("db is nil")
	}

	q := squirrel.
		Select(
			"a.property_id",
			"a.date",
			"quantity",
			"arrival_allowed",
			"departure_allowed",
			"a.minimum_stay as a_min_stay",
			"a.maximum_stay as a_max_stay",
			"COALESCE(duration, 0) as duration",
			"COALESCE(amount, 0) as amount",
			"COALESCE(currency, '') as currency",
			"COALESCE(persons, '') as persons",
			"COALESCE(weekdays, '') as weekdays",
			"COALESCE(p.minimum_stay, 0) as p_min_stay",
			"COALESCE(p.maximum_stay, 0) as p_max_stay",
			"COALESCE(extra_person_price, 0) as extra_person_price",
			"COALESCE(extra_person_price_currency, '') as extra_person_price_currency",
		).
		From(fmt.Sprintf("availabilities as a")).
		RightJoin("prices as p on (p.property_id = a.property_id AND a.date BETWEEN p.period_from AND p.period_till)").
		OrderBy("a.date", "persons", "duration", "p.minimum_stay")

	if len(f.PropertyIDs) > 0 {
		q = q.Where(squirrel.Eq{"property_id": f.PropertyIDs})
	}

	if !f.DateFrom.IsZero() {
		q.Where(squirrel.GtOrEq{"date": f.DateFrom})
	}

	if !f.DateTill.IsZero() {
		q.Where(squirrel.Lt{"date": f.DateTill})
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

	var date string
	result := make([]Model, 0, defaultSize)
	for rows.Next() {
		var mdl Model
		if err := rows.Scan(
			&mdl.PropertyID,
			&date,
			&mdl.Quantity,
			&mdl.ArrivalAllowed,
			&mdl.DepartureAllowed,
			&mdl.ArrivalStayMin,
			&mdl.ArrivalStayMax,
			&mdl.Duration,
			&mdl.Amount,
			&mdl.Currency,
			&mdl.Persons,
			&mdl.Weekdays,
			&mdl.PriceStayMin,
			&mdl.PriceStayMax,
			&mdl.XtraPersonPrice,
			&mdl.XtraPersonPriceCur,
		); err != nil {
			return nil, err
		}

		mdl.Date, err = time.Parse(timeLayoutMysql, date)
		if err != nil {
			return nil, err
		}

		result = append(result, mdl)
	}

	return result, nil
}
