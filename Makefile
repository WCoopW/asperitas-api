include .env
export

DATABASE_URL=postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable

.PHONY: migrate-up migrate-down

migrate-up:
	goose -dir ./internal/db/migrations postgres ${DATABASE_URL} up

migrate-down:
	goose -dir ./internal/db/migrations postgres ${DATABASE_URL} down