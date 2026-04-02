# ---------- 1. BUILD STAGE ----------
FROM golang:1.25.8 AS Builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

# ---------- 2. RUNTIME STAGE ----------
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

ENTRYPOINT ["./app"]