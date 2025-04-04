build:
	@go build -o bin/app ./*.go

run: build
	@./bin/app

test:
	@go test -v ./...