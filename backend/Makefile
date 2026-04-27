# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
.PHONY: audit
audit:
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)"
	go vet ./...

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

## build: build the application
.PHONY: build
build:
	go build -o=/tmp/bin/server ./cmd/

## run: build and run the application
.PHONY: run
run: build
	/tmp/bin/server

## dev: run the application with live reloading
.PHONY: dev
dev:
	air
