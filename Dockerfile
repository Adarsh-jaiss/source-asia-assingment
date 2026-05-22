FROM golang:1.26.1-alpine AS base

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./

RUN go mod download

FROM base AS build

COPY . .

# Build the application using the root main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o app ./main.go

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache make && adduser -D appuser && mkdir -p /app/logs && chown -R appuser:appuser /app

COPY --from=build /app/app .
COPY --from=build /app/makefile ./Makefile
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER appuser

EXPOSE 8080

CMD ["./app"]