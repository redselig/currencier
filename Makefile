.EXPORT_ALL_VARIABLES:

COMPOSE_CONVERT_WINDOWS_PATHS=1

PG_SCHEMAPATH = ./internal/data/repository/db/sql_dev

tidy:
	go mod tidy
fmt:
	go fmt ./...
prepare_lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0
lint: prepare_lint fmt tidy
	golangci-lint run ./...
run:
	go run main.go start
build:
	go build -o currencier.exe main.go
test:
	go test -race ./...
testv:
	go test -v -race ./...
docker-up:
	docker-compose -f docker-compose.yml up -d
docker-down:
	docker-compose -f docker-compose.yml down
docker-down-hard:
	docker-compose -f docker-compose.yml down -v --remove-orphans
	docker rmi currencier:latest