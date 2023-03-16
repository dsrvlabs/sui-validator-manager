.PHONY: build

build: 
	@go build

test:
	@go test ./... -v

coverage:
	@go test ./... -test.coverprofile cover.out

clean:
	@go clean
