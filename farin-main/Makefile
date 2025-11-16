include .env
MIGRATE=migrate -path=migrations -database "$(DATABASE_HOST)" -verbose

devtools:
	@echo "Installing devtools"
	go install github.com/swaggo/swag/cmd/swag@latest
	go get github.com/google/wire/cmd/wire@latest

generate:
	go mod tidy
	go generate ./...
	go mod tidy

db-migrate-up:
		$(MIGRATE) up
db-migrate-down:
		$(MIGRATE) down
db-force:
		@read -p  "Which version do you want to force?" VERSION; \
		$(MIGRATE) force $$VERSION

db-goto:
		@read -p  "Which version do you want to migrate?" VERSION; \
		$(MIGRATE) goto $$VERSION

db-drop:
		$(MIGRATE) drop

db-create-migration:
		@read -p  "What is the name of migration?" NAME; \
		${MIGRATE} create -ext sql -seq -dir migrations  $$NAME

db-seed:
	go run ./cmd/.  --seed --fakeData


swagger:
	swag init --parseDependency -g ./cmd/boot.go -o ./docs