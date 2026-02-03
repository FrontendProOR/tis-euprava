# ===== BUILD STAGE =====
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Kopiraj go mod fajlove
COPY go.mod ./ 
RUN go mod download

# Kopiraj ceo servis
COPY . .

# Build binarni fajl
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o mup-gradjani ./cmd

# ===== RUNTIME STAGE =====
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/mup-gradjani .

EXPOSE 8002

CMD ["./mup-gradjani"]
