.EXPORT_ALL_VARIABLES:

COMPOSE_CONVERT_WINDOWS_PATHS=1

PG_SCHEMAPATH = ./internal/data/repository/db/sql
PGCONNSTRING = "host=localhost port=5432 user=postgres password=postgres dbname=currencier sslmode=disable"
MIGRATIONPATH = ./internal/data/repository/db/migrations

timeout:
	timeout 5
tidy:
	go mod tidy
fmt:
	go fmt ./...
prepare_lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0
prepare_migrate:
	go get -u github.com/pressly/goose/cmd/goose@v2.7.0-rc4
lint: prepare_lint fmt tidy
	golangci-lint run ./...
run:
	go run main.go --debug start
build:
	go build -o currencier.exe main.go
test:
	go test -race ./...
testv:
	go test -v -race ./...

migrate_db:timeout goose_up
goose_up:
	goose -dir $(MIGRATIONPATH) postgres $(PGCONNSTRING) up
goose_down:
	goose -dir $(MIGRATIONPATH) postgres $(PGCONNSTRING) down
docker-up:
	docker-compose -f docker-compose.yml up -d
docker-down:
	docker-compose -f docker-compose.yml down
docker-down-hard:
	docker-compose -f docker-compose.yml down -v --remove-orphans
	docker rmi currencier_currencier:latest
up: docker-up prepare_migrate migrate_db
down: docker-down-hard