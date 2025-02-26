# Tatu

Secret manager for dev teams.

## Development

This section provides guidance for self-hosting and contributing to the project.
Thank you for supporting `tatu`'s development ❤️.

### Requirements

1. Install the dependencies.

```sh
go mod tidy
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
```
### Gotchas

If using `go run main.go ...` for testing the CLI, be aware that if an error occurs, there will be a
`exit status 1` displayed. This is not the case for the actual binary, so don't bother fixing it.
