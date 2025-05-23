include .env
export $(shell sed 's/=.*//' .env)

.PHONY: migrate-up
migrate-up:
	@echo "Running DB migration on: $(DB_ADDR)"
	migrate -database "$(DB_ADDR)" -path cmd/migrate/migrations up


.PHONY: migrate-down
migrate-down:
	@echo "Running DB migration on: $(DB_ADDR)"
	migrate -database "$(DB_ADDR)" -path cmd/migrate/migrations down

.PHONY: swag init
swag init:
	@echo "Generating Swagger docs..."
	swag init --dir cmd/api/ --parseDependency --parseInternal --parseDepth 3
	@echo "Swagger docs generated successfully."