include .env
export $(shell sed 's/=.*//' .env)

.PHONY: migrate-up

migrate-up:
	@echo "Running DB migration on: $(DB_URL)"
	migrate -database "$(DB_URL)" -path ./migrations up
