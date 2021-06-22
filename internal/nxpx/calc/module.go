package calc

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"nxpx/internal/pkg/calc"
	"nxpx/internal/pkg/repo/aprepo"
	"nxpx/internal/pkg/storage"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			func(s *storage.Storage, l *zap.Logger) *aprepo.Repo {
				return aprepo.New(s, l)
			},
		),
		fx.Provide(
			func(r *aprepo.Repo) *calc.Calculator {
				return calc.New(r)
			},
		),
		fx.Invoke(
			func(c *calc.Calculator, l *zap.Logger) error {
				tbl, err := c.Calculate(
					context.Background(),
					time.Date(2017, 5, 18, 0, 0, 0, 0, time.UTC),
				)
				if err != nil {
					return err
				}

				l.Info(fmt.Sprintf("%d loaded", len(tbl.Rows)))
				return nil
			},
		),
	)
}
