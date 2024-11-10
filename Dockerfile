FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

COPY . .
RUN go build -o main ./cmd/kafctl

# Stage 2: Final Image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main ./

COPY app_config.json ssl_config.json ./

EXPOSE 8989
CMD ["./main"]