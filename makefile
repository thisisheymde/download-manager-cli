build:
	@go build -o bin/ddman

run: build
	@./bin/ddman

test:
	@go test -v ./...