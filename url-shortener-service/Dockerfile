FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener-service ./cmd/server/main.go
RUN chmod +x url-shortener-service

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/url-shortener-service .

EXPOSE 50001

CMD ["./url-shortener-service"]