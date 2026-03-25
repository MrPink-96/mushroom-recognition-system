# ---------- 1. BUILD STAGE ----------
FROM golang:1.24.7 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go bild -o app ./cmd/main.go


# ---------- 2. RUNTIME STAGE ----------
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8081

CMD ["./app"]