include development.env

DB_DSN = ${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}
MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: services-up
services-up:
	docker compose -f docker-compose.yaml --env-file ./development.env up --build

.PHONY: services-down
services-down:
	docker compose -f docker-compose.yaml --env-file ./development.env down -v

.PHONY: migration
migration:
	@migrate create -seq -ext sql -dir ${MIGRATIONS_PATH} ${filter-out $@,${MAKECMDGOALS}}

.PHONY: migrate-up
migrate-up:
	@migrate -path ${MIGRATIONS_PATH} -database ${DB_DSN} -verbose up

.PHONY: migrate-down
migrate-down:
	@migrate -path ${MIGRATIONS_PATH} -database ${DB_DSN} -verbose down

.PHONY: seed
seed:
	@go run cmd/migrate/seed/main.go

PHONY: test
test:
	@go test -v ./...

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt
