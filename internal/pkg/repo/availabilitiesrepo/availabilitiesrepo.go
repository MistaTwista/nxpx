package availabilitiesrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

const (
	defaultSize    = 64
	defaultTimeout = 1 * time.Second
)

type Model struct {
	ID               int
	PropertyID       string
	Date             time.Time
	Quantity         int
	ArrivalAllowed   bool
	DepartureAllowed bool
	StayMin          int
	StayMax          int
	Version          uint
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
		log:       log.Named("availabilities-repo"),
	}
}

type Filter struct {
	PropertyIDs []string
	DateFrom    time.Time
	DateTill    time.Time
}

func (r *Repo) List(ctx context.Context, f Filter) ([]Model, error) {
	q := squirrel.
		Select("id",
			"property_id",
			"date",
			"quantity",
			"arrival_allowed",
			"departure_allowed",
			"minimum_stay",
			"maximum_stay",
			"version",
		).
		From(r.tableName)

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

	result := make([]Model, 0, defaultSize)
	for rows.Next() {
		var mdl Model
		if err := rows.Scan(
			&mdl.ID,
			&mdl.PropertyID,
			&mdl.Quantity,
			&mdl.ArrivalAllowed,
			&mdl.DepartureAllowed,
			&mdl.StayMin,
			&mdl.StayMax,
			&mdl.Version,
		); err != nil {
			return nil, err
		}

		result = append(result, mdl)
	}

	return result, nil
}
