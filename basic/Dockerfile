# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/ .
RUN go build -o main .

# Stage 2: Run the Go binary in a lightweight container
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
