# Build

FROM golang:1.22.0-alpine3.18 as BUILD

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/prospector

# Binaries

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/bin/prospector /app/prospector

CMD ["./prospector", "server"]
