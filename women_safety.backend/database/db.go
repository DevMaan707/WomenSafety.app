package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB(dbURL string) error {
	var err error
	DB, err = sqlx.Connect("postgres", dbURL)
	if err != nil {
		return err
	}
	return nil
}
