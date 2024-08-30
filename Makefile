create-migration:
	@echo "Creating migration..."
	migrate create -ext sql -dir internal/db/migrations -seq -digits 1 $(name)

# Run migrations
migrate-up:
	migrate -path internal/db/migrations -database $(DB_URL) up

.PHONY: create-migration
