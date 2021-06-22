package storage

import (
	"nxpx/internal/nxpx/config"
	"nxpx/internal/pkg/storage"

	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			func(cfg *config.Config) (*storage.Storage, error) {
				return storage.New(cfg.Storage)
			},
		),
		fx.Invoke(
			func(lf fx.Lifecycle, s *storage.Storage) {
				lf.Append(fx.Hook{
					OnStart: s.Start,
					OnStop:  s.Stop,
				})
			},
		),
	)
}
