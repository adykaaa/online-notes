include development.env
export

help: # Show help for each of the Makefile recipes.
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done
.PHONY: help

pg-up: # Sets up PostgreSQL in a docker container
	docker run -e POSTGRES_USER=user -e POSTGRES_PASSWORD=pgpass123 -e POSTGRES_DB=notes -p 5432:5432 -d --name postgres-dev postgres
.PHONY: pg-up

sqlc: # Generates the DB backend code according to the sqlc.yaml file located in the root folder
	sqlc generate
.PHONY: sqlc

test: # Runs the unit tests
	go test -v -cover ./...
.PHONY: test

mock: # Generates the DB mocks
	mockgen -package mockdb -destination db/mock/querier.go  github.com/adykaaa/online-notes/db/sqlc Querier
.PHONY: mock

build-backend: # Builds the backend Docker image
	docker build . -t online-notes-backend
.PHONY: build-backend

run-backend: # Runs the backend Docker image
	docker run -d -p 8080:8080 online-notes-backend
.PHONY: build-backend

build-frontend: # Builds the frontend Docker image
	docker build ./web/ -t online-notes-frontend
.PHONY: build-frontend
