BINARY_NAME=prospector

.PHONY: docs vet test test_coverage clean build

build:
	CGO_ENABLED=0 GOOS=linux go build -o bin/prospector

clean:
	go clean
	rm ./bin/${BINARY_NAME}

test:
	go test -v ./... -race

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

cov_percent: coverage
	go tool cover -func coverage.out | sed -n -e '/^total/s/:.*statements)[^0-9]*/: /p'

docs:
	swag fmt && swag init -g routes/routes.go
