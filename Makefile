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
	@docker compose -f compose.prettier.yml run prettier
.PHONY: prettier

test:
	@GO_ENV=test go test -count=1 -v  ./test/...
.PHONY: test

lint-fix:
	@golangci-lint run -c .golangci-gin.yml --fix
.PHONY: lint

