# Build

FROM golang:1.21.3-alpine3.17 as BUILD

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/prospector

# Binaries

FROM alpine:latest

COPY --from=build /app/bin/prospector /app/bin/prospector

CMD ["/app/bin/prospector", "server"]