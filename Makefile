include development.env
export

help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

pg-up:
	docker run -e POSTGRES_USER=user -e POSTGRES_PASSWORD=pgpass123 -e POSTGRES_DB=notes -p 5432:5432 -d --name postgres-dev postgres
.PHONY: pg-up

sqlc:
	sqlc generate
.PHONY: sqlc

test:
	go test -v -cover ./...
.PHONY: test

mock:
	mockgen -package mockdb -destination db/mock/querier.go  github.com/adykaaa/online-notes/db/sqlc Querier
.PHONY: mock