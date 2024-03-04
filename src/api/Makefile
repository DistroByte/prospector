BINARY_NAME=prospector

.PHONY: docs vet test test_coverage clean build

build:
	CGO_ENABLED=0 GOOS=linux go build -o bin/prospector

clean:
	go clean
	rm ./bin/${BINARY_NAME}

test:
	go test -v ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

docs:
	swag fmt && swag init -g routes/routes.go
