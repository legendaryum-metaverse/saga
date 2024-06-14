# Makefile

all: format lint-fix
.PHONY: all

format:
	@gofmt -l -s -w .
.PHONY: format

lint:
	@golangci-lint run -c .golangci-gin.yml
.PHONY: lint

test:
	@GO_ENV=test go test -count=1  ./test/...
.PHONY: test

lint-fix:
	@golangci-lint run -c .golangci-gin.yml --fix
.PHONY: lint

