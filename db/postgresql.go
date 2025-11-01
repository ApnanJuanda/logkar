package db

import (
	"bsnack/lib/env"
	"database/sql"
	"log"
)

func PostgresqlOpen() *sql.DB {
	var sqlDb *sql.DB
	var err error
	sqlDb, err = sql.Open("postgres", env.String("POSTGRESQL_URL", "postgresql://postgres:postgres@127.0.0.1:5432/db_account?sslmode=disable"))
	if err != nil {
		log.Fatal(err)
	}
	return sqlDb
}
