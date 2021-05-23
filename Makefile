BINPATH       = $(PWD)/bin

.PHONY: deps
deps: deps-tools deps-main

.PHONY: deps-tools
deps-tools:
	cd tools && go mod tidy && go mod vendor && go generate tools.go

.PHONY: deps-main
deps-main:
	go mod tidy && go mod vendor

.PHONY: build
build:
	go build -o $(BINPATH)/procrastination-killer main.go

.PHONY: lint
lint:
	${BINPATH}/golangci-lint run -v