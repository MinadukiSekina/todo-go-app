.PHONY: build

build:
	go mod tidy

gotest:
	@make build
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o cover.html