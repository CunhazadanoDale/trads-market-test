DB_URL=postgres://trads:trads@postgres:5432/trads?sslmode=disable

.PHONY: wait-for-db migrate-up migrate-down migrate-status import-ans

wait-for-db:
	docker-compose up -d postgres
	@echo "Waiting for Postgres to be ready..."
	@until docker-compose exec -T postgres pg_isready -U trads -d trads >/dev/null 2>&1; do \
	  printf "."; \
	  sleep 1; \
	done; \
	echo " Postgres is ready."

migrate-up: wait-for-db
	docker-compose run --rm migrate -path=/migrations -database "$(DB_URL)" up

migrate-down: wait-for-db
	docker-compose run --rm migrate -path=/migrations -database "$(DB_URL)" down

migrate-status: wait-for-db
	docker-compose run --rm migrate -path=/migrations -database "$(DB_URL)" version

import-ans:
	docker-compose --profile ans run --rm import_ans
