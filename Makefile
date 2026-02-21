SERVICE_NAME := main

init:
	make clean

tidy:
	go mod tidy

vendor:
	make clean
	go mod tidy
	go mod vendor

clean-dbs:
	rm -rf dbs

clean-build:
	rm -f $(SERVICE_NAME)

clean:
	make clean-dbs
	make clean-build
	rm -rf vendor

start:
	air

run-migration-up:
	go run cmd/main.go --migrate

new-migration:
	migrate create -ext sql -dir internal/config/db/migrations -seq $(NAME)

swag:
	swag init -d ./cmd,./internal/application/rest/controllers,./internal/application/dto,internal/models -o swagger
