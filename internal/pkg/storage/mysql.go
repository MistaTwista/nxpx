package storage

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	ConnString string `envconfig:"db_conn_string" default:"user:password@tcp(localhost:3306)/db"`
}

func New(c Config) *Storage {
	return &Storage{
		config: c,
	}
}

type Storage struct {
	config Config
	db     *sql.DB
}

func (s *Storage) GetDB() *sql.DB {
	return s.db
}

func (s *Storage) Start(_ context.Context) error {
	db, err := sql.Open("mysql", s.config.ConnString)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Storage) Stop(_ context.Context) error {
	return s.db.Close()

}

func (s *Storage) IsReady() bool {
	err := s.db.Ping()
	if err != nil {
		return false
	}

	return true
}
