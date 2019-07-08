MIGRATION_PATH := ./migrations
DATABASE_URL ?= "postgres://user:password@db:5432/urlshortener?sslmode=disable"

.PHONY: hello
hello:
	@echo "hello"

# Creates a new database migration with a title.
.PHONY: new_migration/%
new_migration/%:
	migrate create -dir $(MIGRATION_PATH) -ext pgsql $*

# Run database migrations.
# To run all migrations use `latest`.
.PHONY: migrate
migrate:
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) up

# Run an assembly in `cmd`
.PHONY: run/%
run/%:
	go run ./cmd/$*/*.go $(EXTRA_ARGS)

