package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"nxpx/internal/nxpx/config"
	"nxpx/internal/pkg/calc"
	"nxpx/internal/pkg/logger"
	"nxpx/internal/pkg/repo/aprepo"
	"nxpx/internal/pkg/storage"
)

var (
	// Provisioned by ldflags
	version   string
	buildDate string
	commit    string
)

const appName = "nxpx"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--help" {
		if err := config.Usage(); err != nil {
			log.Println(err.Error())
		}

		return
	}

	conf, err := config.New()
	if err != nil {
		panic(err)
	}

	conf.Version = version
	conf.BuildDate = buildDate
	conf.Commit = commit

	zapLogger, err := logger.New(conf.Logger, appName, version, buildDate, commit)
	if err != nil {
		panic(err)
	}

	db := storage.New(conf.Storage)
	err = db.Start(context.Background())
	if err != nil {
		panic(err)
	}

	repo := aprepo.New(db, zapLogger)

	c := calc.New(repo)
	table, err := c.Calculate(
		context.Background(),
		time.Date(2017, 5, 18, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		panic(err)
	}

	zapLogger.Info(fmt.Sprintf("%d loaded", len(table.Rows)))
}
