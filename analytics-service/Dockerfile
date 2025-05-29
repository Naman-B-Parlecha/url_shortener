FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o analytics-service ./cmd/server/main.go
RUN chmod +x analytics-service

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/analytics-service .

EXPOSE 50002

CMD ["./analytics-service"]