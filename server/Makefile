.PHONY: setup-db stop-db delete-db reboot-db

build:
	@go build -o server ./main.go

run-dev:
	@go run main.go

install-deps:
	@go mod download

setup-db:
	@docker run --name postgres -e POSTGRES_PASSWORD=pwd -p 5432:5432 -itd postgres:latest

stop-db:
	@docker stop postgres

delete-db:
	@docker rm postgres

reboot-db: stop-db delete-db setup-db

migrate:
	@docker cp ./schema.sql postgres:/tmp/schema.sql
	@docker exec -it postgres psql -U postgres -c "CREATE DATABASE visio;"
	@docker exec -it postgres psql -U postgres -d visio -f /tmp/schema.sql
	@echo "Database rebooted"

test:
	@go test -v ./...
