MIGRATIONS_PATH=migrations/postgres
DB_URL=postgresql://admin:admin@localhost:5432/analysis_db?sslmode=disable

services-up:
	docker-compose up

services-down:
	docker-compose down

migrate-up:
	goose -dir=$(MIGRATIONS_PATH) postgres $(DB_URL) up

migrate-down:
	goose -dir=$(MIGRATIONS_PATH) postgres $(DB_URL) down