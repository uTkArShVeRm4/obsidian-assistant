FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

RUN apk add --no-cache git tzdata

WORKDIR /app

COPY --from=builder /app/main .
COPY .env .

# Create directories
RUN mkdir -p /app/data/uploads /app/data/Main/Attachments

EXPOSE 7777

# We'll use an entrypoint script to configure git
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
