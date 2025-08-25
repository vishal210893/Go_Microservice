# Variables
MIGRATIONS_PATH = ./cmd/migrate/migrations
DATABASE_URL = postgres://avnadmin:AVNS_LT5DsEKUPKfrHSHZHyB@pg-1d9d15dc-vishal210893-5985.h.aivencloud.com:28832/defaultdb?sslmode=require

# Migration commands
.PHONY: migrate-create-users migrate-create-posts migrate-up migrate-down migrate-force migrate-version seed gen-docs

migrate-create-users:
	migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) create_users

migrate-create-posts:
	migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) create_posts

migrate-up:
	migrate -path=$(MIGRATIONS_PATH) -database="$(DATABASE_URL)" up

migrate-down:
	migrate -path=$(MIGRATIONS_PATH) -database="$(DATABASE_URL)" down

migrate-force:
	migrate -path=$(MIGRATIONS_PATH) -database="$(DATABASE_URL)" force $(version)

migrate-version:
	migrate -path=$(MIGRATIONS_PATH) -database="$(DATABASE_URL)" version

# Generic migration creation (usage: make migrate-create name=table_name)
migrate-create:
	migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(name)

seed:
	go run cmd/migrate/seed/main.go

gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt

# Help command
help:
	@echo "Available commands:"
	@echo "  migrate-create-users  - Create users migration files"
	@echo "  migrate-create-posts  - Create posts migration files"
	@echo "  migrate-create name=X - Create migration with custom name"
	@echo "  migrate-up           - Apply all up migrations"
	@echo "  migrate-down         - Apply all down migrations"
	@echo "  migrate-force        - Force migration version (usage: make migrate-force version=N)"
	@echo "  migrate-version      - Show current migration version"