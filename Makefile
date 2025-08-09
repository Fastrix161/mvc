build:
	@go build -o bin/MVC cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/MVC