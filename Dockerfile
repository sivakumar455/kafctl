FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk update
RUN apk add \
    gcc \
    # use mold for convenient extra linker inputs
    mold \
    musl-dev \
    # explicitly install SASL package
    cyrus-sasl-dev

COPY ./cmd /app/cmd
COPY ./internal /app/internal
COPY ./web /app/web
COPY go.mod go.sum ./

COPY app_config.json ssl_config.json ./

RUN go mod download

RUN go build -tags musl -o main ./cmd/kafctl

# Stage 2: Final Image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ ./

COPY app_config.json ssl_config.json ./

EXPOSE 8989
CMD ["./main"]