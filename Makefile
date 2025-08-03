run:
	go run cmd/main.go

generate:
	sqlc generate -f internal/db/sqlc.yaml

.PHONY: migration

migration:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir internal/db/migrations -seq "$$(date +%Y%m%d_%H%M%S)_$${name}"

generate-swagger:
	swagger init --out docs --generateInfo cmd/main.go