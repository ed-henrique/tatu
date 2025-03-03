/*
Copyright © 2025 Eduardo Henrique Freire Machado

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package db

import (
	"database/sql"
	"embed"
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var embeddedMigrations embed.FS

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

	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(embeddedMigrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		fmt.Fprintf(os.Stderr, "Could not set dialect: %s", err.Error())
		os.Exit(1)
		return nil
	}

	if err := goose.Up(dbConn, "migrations"); err != nil {
		fmt.Fprintf(os.Stderr, "Could not migrate: %s", err.Error())
		os.Exit(1)
		return nil
	}

	return dbConn
}
