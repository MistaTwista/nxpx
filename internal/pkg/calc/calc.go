package calc

import (
	"context"
	"sort"
	"time"

	"nxpx/internal/pkg/cal"
	"nxpx/internal/pkg/repo/aprepo"
)

type apRepo interface {
	List(ctx context.Context, f aprepo.Filter) ([]aprepo.Model, error)
}

const (
	calcDays = 365
	oneDay   = 24 * time.Hour
	calcLen  = calcDays * oneDay
)

type Calculator struct {
	apRepo apRepo
	cal    *cal.Calendar
}

func New(apRepo apRepo) *Calculator {
	return &Calculator{
		apRepo: apRepo,
	}
}

//"71438849-47cb-4b00-82de-34fff691f017"
func (c *Calculator) Calculate(ctx context.Context, t time.Time, propID ...string) (Table, error) {
	calcTo := t.Add(calcLen)
	prices, err := c.apRepo.List(ctx, aprepo.Filter{
		PropertyIDs: propID,
		DateFrom:    t,
		DateTill:    calcTo,
	})
	if err != nil {
		return Table{}, err
	}

	return MakeMeHappy(cal.New(prices), t, calcTo, calcDays), nil
}

func beginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

type calendar interface {
	At(time.Time) (cal.Date, bool)
}

func MakeMeHappy(cal calendar, from, to time.Time, calculationDays int) Table {
	rows := make([]Row, 0, calculationDays)
	now := beginningOfDay(from)

	for now.Before(beginningOfDay(to)) {
		iterDate := now
		now = now.Add(oneDay)
		d, ok := cal.At(iterDate)
		if !ok {
			rows = append(rows, Row{
				ArrivalDate:  iterDate,
				PersonsCount: 1,
				Stays:        make([]Stay, calculationDays),
			})
			continue
		}

		// TODO: look next day or prev day for some rules
		// TODO: get durations (how to
		for _, price := range d.Prices {
			for _, personsCount := range price.ByPerson() {
				stays := make([]Stay, 0, calculationDays)
				for i := 0; i < calculationDays; i++ {
					currentDate := iterDate.Add(time.Duration(i) * 24 * time.Hour)
					_, ok := cal.At(currentDate)
					if !ok {
						stays = append(stays, Stay{
							Date:  currentDate,
							Price: 0,
						})
						continue
					}

					// d.priceFor(personsCount)
					stays = append(stays, Stay{
						Date:  currentDate,
						Price: price.Amount + price.XtraPersonPrice*(personsCount-1),
					})
				}

				rows = append(rows, Row{
					ArrivalDate:  iterDate,
					PersonsCount: uint(personsCount),
					Stays:        stays,
				})
			}
		}
	}

	return Table{Rows: rows}
}

/*
- get data from db that:
today + 365 days, duration 1 - 21 days
we can ask person count OR 1...max (6 in DB)
we can ask property ID
- DB is slow, so take data few times as possible
- no price (not in DB) OR not valid price (bad min_stay) - show 0 price
- table starts always 1...N, show 0 prices where no data
- if arrival_allowed (= 0) for date any durations started from it - have 0 price
... TODO: add other rules here
*/

type Stay struct {
	Date time.Time
	//Price decimal.Decimal
	Price int
}

type Row struct {
	ArrivalDate  time.Time
	PersonsCount uint
	Stays        []Stay
}

// Table shown for one property at a time
type Table struct {
	Rows []Row
}

// breakDown split number of days into list of available durations
// breakDown([]int{7,3,1}, 9) => []int{7,1,1}
func breakDown(durations []int, days int) []int {
	sort.Ints(durations)
	res := make([]int, 0, len(durations))
oLoop:
	for i := len(durations) - 1; i >= 0; i-- {
		itm := durations[i]

	iLoop:
		for {
			if days <= 0 {
				break iLoop
			}

			if days-itm < 0 {
				continue oLoop
			}

			days -= itm
			res = append(res, itm)
		}
	}

	return res
}

//// Prices index for data to search fast
//type Prices struct {
//	list      []pricerepo.Model
//	durations map[uint]struct{}
//}
//
//func NewPrices(list []pricerepo.Model) *Prices {
//	durs := make(map[uint]struct{})
//	for _, p := range list {
//		durs[uint(p.Duration)] = struct{}{}
//	}
//
//	return &Prices{
//		list:      list,
//		durations: durs,
//	}
//}
