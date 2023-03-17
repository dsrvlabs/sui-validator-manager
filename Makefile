.PHONY: build

all: build test coverage

build: 
	@go build

test:
	@go test ./... -v

coverage:
	@go test ./... -test.coverprofile cover.out

fmt:
	@go fmt ./...

clean:
	@go clean
