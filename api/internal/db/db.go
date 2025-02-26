package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func New(dsn string) *sql.DB {
	dbConn, err := sql.Open("sqlite", dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open DB: %s", err.Error())
		os.Exit(1)
		return nil
	}

	if err := dbConn.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not ping DB: %s", err.Error())
		os.Exit(1)
		return nil
	}

	return dbConn
}
