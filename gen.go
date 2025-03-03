package main

//go:generate goose -dir ./internal/models/migrations sqlite3 ./db.sqlite up
//go:generate sqlc generate -f internal/models/sqlc.yaml
