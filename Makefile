include .envrc

## run/hn: run the cmd/hn application
.PHONY: run/api
run/hn:
	go run ./cmd/hackernews

## run/fetcher: run the cmd/fetcher application
.PHONY: run/fetcher
run/fetcher:
	go run ./cmd/fetcher -db-dsn=${DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${DB_DSN} up

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${DB_DSN}

.PHONY: test
test:
	go test ./...

