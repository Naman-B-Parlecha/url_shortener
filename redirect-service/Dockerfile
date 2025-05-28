FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o redirect-service ./cmd/server/main.go
RUN chmod +x redirect-service

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/redirect-service .

EXPOSE 50002

CMD ["./redirect-service"]