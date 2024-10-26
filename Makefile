ifneq (,$(wildcard ./.env))
    include .env
    export
endif

DB_DSN="postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTRGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"

all: run

migrations-up: run-storage
	go install github.com/pressly/goose/v3/cmd/goose@latest
	goose -dir migrations postgres $(DB_DSN) up

migrations-down: run-storage
	goose -dir migrations postgres $(DB_DSN) down

run:
	docker compose up --build --remove-orphans

stop:
	docker compose down

run-local-storage:
	docker compose -f docker-compose.local.yaml up -d

stop-local-storage:
	docker compose -f docker-compose.local.yaml down

run-local:
	go run cmd/main.go

build:
	go build -o build/main cmd/main.go

clean:
	rm build/*

rebuild: clean build