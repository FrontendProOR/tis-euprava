FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o sso-auth ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/sso-auth .

EXPOSE 8001

CMD ["./sso-auth"]
