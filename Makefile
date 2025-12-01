.PHONY: run build test docker-up docker-down migrate

run:
	go run ./cmd/api api

build:
	go build -o bin/api ./cmd/api
	go build -o bin/migrate ./cmd/migrate

test:
	go test -v -cover ./...

migrate:
	go run ./cmd/migrate

docker-up:
	docker-compose up --build -d
	@echo "Waiting for MySQL..."
	@sleep 10
	docker-compose run --rm migrate
	@echo "Ready! API: http://localhost:8080"

docker-down:
	docker-compose down

logs:
	docker-compose logs -f api
