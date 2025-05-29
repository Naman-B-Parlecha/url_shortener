FROM golang:alpine AS builder

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o broker-service ./cmd/server/main.go
RUN chmod +x broker-service

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache curl

COPY --from=builder /app/broker-service .

EXPOSE 8081

CMD ["./broker-service"]