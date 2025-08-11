run:
	go run cmd/main.go

generate:
	sqlc generate -f internal/db/sqlc.yaml

.PHONY: migration

migration:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir internal/db/migrations -seq "$$(date +%Y%m%d_%H%M%S)_$${name}"

generate-swagger:
	swag init -g cmd/main.go -o docs



include .env
export $(shell sed 's/=.*//' .env)
.PHONY: migrate-up

migrate-up:
	migrate -path internal/db/migrations \
	  -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" \
	  up