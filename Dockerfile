FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o matrix-compute ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/matrix-compute .

EXPOSE 8080

CMD ["./matrix-compute"] 
