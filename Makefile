include .envrc

## run/hn: run the cmd/hn application
.PHONY: run/api
run/hn:
	go run ./cmd/hackernews

## run/fetcher: run the cmd/fetcher application
.PHONY: run/fetcher
run/fetcher:
	go run ./cmd/fetcher

