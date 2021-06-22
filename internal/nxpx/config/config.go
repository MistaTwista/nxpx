package config

import (
	"nxpx/internal/pkg/logger"
	"nxpx/internal/pkg/storage"

	"github.com/kelseyhightower/envconfig"
)

const (
	EnvPrefix = "NXPX"
	EnvDev    = "dev"
)

type Config struct {
	Env       string         `envconfig:"env"`
	Debug     bool           `envconfig:"debug"`
	Logger    logger.Config  `envconfig:"logger"`
	Storage   storage.Config `envconfig:"storage"`
	Version   string
	BuildDate string
	Commit    string
}

func Usage() error {
	return envconfig.Usage(EnvPrefix, &Config{})
}

func (c Config) IsDev() bool {
	return c.Env == EnvDev
}

func (c *Config) Validate() error {
	// TODO: add more validation here

	return nil
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process(EnvPrefix, cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	if cfg.Debug {
		cfg.Logger.Level = "debug"
		cfg.Logger.Debug = true
	}

	return cfg, nil
}
