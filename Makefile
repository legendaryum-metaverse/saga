# Makefile

all: format lint-fix prettier test
.PHONY: all

format:
	@gofmt -l -s -w .
.PHONY: format

lint:
	@golangci-lint run -c .golangci-gin.yml
.PHONY: lint

prettier:
	@bun i && bun run format
.PHONY: prettier

prettier-build:
	@COMPOSE_PROJECT_NAME=golib-prettier docker compose -f compose.prettier.yml --progress=plain build prettier
.PHONY: prettier-build

test:
	@docker compose up -d
	@GO_ENV=test go test -count=1 -v ./test/...
.PHONY: test

lint-fix:
	@golangci-lint run -c .golangci-gin.yml --fix
.PHONY: lint

