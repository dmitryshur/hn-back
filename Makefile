include .envrc

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/hackernews

