# ===== Build stage =====
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/api

# ===== Runtime stage =====
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/app /app/app
COPY db/migrations /app/db/migrations

EXPOSE 8080

CMD ["/app/app"]
