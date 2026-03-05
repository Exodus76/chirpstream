#Makefile
all: build

build-apigateway:
	go build -ldflags="-s -w" -o bin/chirpstream-api-gateway cmd/api-gateway/main.go

build-user-service:
	go build -ldflags="-s -w" -o bin/chirpstream-user cmd/user-service/main.go

build-rel-service:
	go build -ldflags="-s -w" -o bin/chirpstream-relationship cmd/user-relationship/main.go
# 	go build -tags="gocql_debug" -ldflags="-s -w" -o bin/chirpstream-relationship cmd/user-relationship/main.go

build-chirp-service:
	go build -ldflags="-s -w" -o bin/chirpstream-chirps cmd/chirp-service/main.go
# 	go build -tags="gocql_debug" -ldflags="-s -w" -o bin/chirpstream-chirps cmd/chirp-service/main.go

build: build-apigateway build-user-service build-chirp-service

run-apigateway: build-apigateway
	./bin/chirpstream-api-gateway

run-user: build-user-service
	./bin/chirpstream-user

run-rel: build-rel-service
	./bin/chirpstream-relationship

run-chirp: build-chirp-service
	./bin/chirpstream-chirps

SHELL := /bin/bash

test-user-service:
	@docker compose up -d
	@trap 'echo "removing db"; docker compose -f docker-compose.yml down' EXIT
	@go test -v ./internal/user

clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod download

# Install sqlc and goose
install-tools:
	go install github.com/pressly/goose/v3/cmd/goose@latest

# Run goose migrations
migrate:
	goose up

# Rollback the last migration
migrate-down:
	goose down

# Create a new migration file
# eg: make migrate-create name="template_group" 
migrate-create:
	goose -dir ./internal/data/migrations create $(name) sql

#Proto generation
generate-proto:
	protoc --proto_path=proto --go_out=internal/pb --go_opt=paths=source_relative --go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative user.proto


# Help command to list all available targets
# @echo "  test          - Run tests"
help:
	@echo "Available targets:"
	@echo "  all            - Build the application (default target)"
	@echo "  build          - Build the Go application"
	@echo "  run            - Run the Go application"
	@echo "  migrate        - Run goose migrations"
	@echo "  migrate-down   - Rollback the last migration"
	@echo "  migrate-create - Create a new migration file (usage: make migrate-create name=<migration_name>)"
	@echo "  clean          - Clean up build files"
	@echo "  deps           - Install Go dependencies"
	@echo "  install-tools  - Install sqlc and goose"
	@echo "  help           - Show this help message"

.PHONY: all build run migrate migrate-down migrate-create clean deps install-tools help
