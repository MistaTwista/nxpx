package cal

import (
	"errors"
	"time"

	"nxpx/internal/pkg/repo/aprepo"
	"nxpx/internal/pkg/repo/availabilitiesrepo"
	"nxpx/internal/pkg/repo/pricerepo"
)

// TODO: make effective search by list of persons
type Date struct {
	Availability availabilitiesrepo.Model
	Prices       []pricerepo.Model
}

//func (d *Date) ByPersonCount() pricerepo.Model {
//}

var ErrNotFound = errors.New("not found")

const (
	oneDay = 24 * time.Hour
)

type Duration struct {
}

type Calendar struct {
	data map[time.Time]*Date
}

func New(list []aprepo.Model) *Calendar {
	data := make(map[time.Time]*Date)
	for _, l := range list {
		bod := beginningOfDay(l.Date)
		if _, ok := data[bod]; !ok {
			data[bod] = &Date{
				Availability: availabilitiesrepo.Model{
					PropertyID:       l.PropertyID,
					Date:             l.Date,
					Quantity:         l.Quantity,
					ArrivalAllowed:   l.ArrivalAllowed,
					DepartureAllowed: l.DepartureAllowed,
					StayMin:          l.ArrivalStayMin,
					StayMax:          l.ArrivalStayMax,
				},
				Prices: nil,
				//Duration: nil,
			}
		}

		data[bod].Prices = append(data[bod].Prices, pricerepo.Model{
			PropertyID:         l.PropertyID,
			Duration:           l.Duration,
			Amount:             l.Amount,
			Currency:           l.Currency,
			Persons:            l.Persons,
			Weekdays:           l.Weekdays,
			StayMin:            l.PriceStayMin,
			StayMax:            l.PriceStayMax,
			XtraPersonPrice:    l.XtraPersonPrice,
			XtraPersonPriceCur: l.XtraPersonPriceCur,
		})
	}

	return &Calendar{data: data}
}

func (c *Calendar) At(t time.Time) (Date, bool) {
	bod := beginningOfDay(t)
	d, ok := c.data[bod]
	if !ok {
		return Date{}, ok
	}

	return *d, ok
}

//func (c *Calendar) DayBefore(t time.Time) (Date, bool) {
//	d, ok := c.data[beginningOfDay(t).Add(-oneDay)]
//	if !ok {
//		return Date{}, ok
//	}
//
//	return *d, ok
//}
//
//func (c *Calendar) DayAfter(t time.Time) (Date, bool) {
//	d, ok := c.data[beginningOfDay(t).Add(oneDay)]
//	if !ok {
//		return Date{}, ok
//	}
//
//	return *d, ok
//}

func beginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
