.PHONY: generate dev

generate:
	go run github.com/99designs/gqlgen generate

dev:
	go run server.go

.DEFAULT_GOAL := dev