.PHONY: build
build:
	go build -v ./cmd/app/vh

.PHONY: migrate
migrate:
	migrate -path migrations -database "postgres://postgres:mg0208@localhost:5432/vh?sslmode=disable" up

.PHONY: create_migration
create_migration:
	migrate create -ext sql -dir migrations -seq -digits 6 $(MG_NAME)

.DEFAULT_GOAL: build