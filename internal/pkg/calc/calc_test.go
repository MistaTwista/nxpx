package calc

import (
	"context"
	"nxpx/internal/pkg/repo/aprepo"
	"nxpx/internal/pkg/storage"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestMakeMeHappy(t *testing.T) {
	t.Skip()
	db := storage.New(storage.Config{ConnString: "user:password@tcp(localhost:3306)/db"})
	err := db.Start(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	repo := aprepo.New(db, &zap.Logger{})

	c := New(repo)
	table, err := c.Calculate(
		context.Background(),
		time.Date(2017, 5, 18, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%d loaded", len(table.Rows))
}

func TestBreakDown(t *testing.T) {
	cases := []struct {
		Name      string
		Durations []int
		Days      int
		Result    []int
	}{
		{
			Name:      "no one",
			Durations: []int{7, 2},
			Days:      5,
			Result:    []int{2, 2},
		},
		{
			Name:      "just5",
			Durations: []int{7, 3, 1},
			Days:      5,
			Result:    []int{3, 1, 1},
		},
		{
			Name:      "just9",
			Durations: []int{7, 3, 1},
			Days:      9,
			Result:    []int{7, 1, 1},
		},
		{
			Name:      "just12",
			Durations: []int{7, 3, 1},
			Days:      12,
			Result:    []int{7, 3, 1, 1},
		},
		{
			Name:      "just12with2",
			Durations: []int{7, 3, 2, 1},
			Days:      12,
			Result:    []int{7, 3, 2},
		},
		{
			Name:      "just13",
			Durations: []int{7, 3, 2, 1},
			Days:      13,
			Result:    []int{7, 3, 3},
		},
		{
			Name:      "just14",
			Durations: []int{7, 3, 2, 1},
			Days:      14,
			Result:    []int{7, 7},
		},
		{
			Name:      "just15",
			Durations: []int{7, 3, 2, 1},
			Days:      15,
			Result:    []int{7, 7, 1},
		},
	}

	for _, c := range cases {
		res := breakDown(c.Durations, c.Days)

		if len(res) != len(c.Result) {
			t.Error("len not equal")
		}
		for i, r := range res {
			if r != c.Result[i] {
				t.Fatalf("%s: want %d, got %d (%v vs %v)", c.Name, r, c.Result[i], res, c.Result)
			}
		}
	}
}
