FROM golang:1.22.7 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/api ./cmd/api

FROM debian:bullseye-slim

ENV SERVER_ADDRESS=0.0.0.0:8080
ENV POSTGRES_USERNAME=postgres
ENV POSTGRES_PASSWORD=1234
ENV POSTGRES_HOST=DB
ENV POSTGRES_PORT=5432
ENV POSTGRES_DATABASE=merch_shop_service
ENV POSTGRES_CONN=postgres://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable
ENV ENV=prod
ENV JWT_SECRET=avito

WORKDIR /app

COPY --from=builder /app/api .

COPY internal/storage/migrations /app/internal/storage/migrations

CMD ["/app/api"]