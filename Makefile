DB_URL=postgres://trads:trads@postgres:5432/trads?sslmode=disable

.PHONY: up down reset wait-for-db migrate-up migrate-down migrate-status import import-ans test lint

up:
	docker compose up -d --build

down:
	docker compose down

reset:
	docker compose down -v

wait-for-db:
	docker compose up -d postgres
	@echo "Aguardando o Postgres ficar pronto..."
	@until docker compose exec -T postgres pg_isready -U trads -d trads >/dev/null 2>&1; do \
	  printf "."; \
	  sleep 1; \
	done; \
	echo " Postgres pronto."

migrate-up: wait-for-db
	docker compose run --rm migrate -path=/migrations -database "$(DB_URL)" up

migrate-down: wait-for-db
	docker compose run --rm migrate -path=/migrations -database "$(DB_URL)" down

migrate-status: wait-for-db
	docker compose run --rm migrate -path=/migrations -database "$(DB_URL)" version

import:
	docker compose run --rm import

import-ans:
	docker compose --profile ans run --rm import_ans

test:
	cd backend-golang && go vet ./... && go test ./...

lint:
	cd frontend && npx oxlint src
