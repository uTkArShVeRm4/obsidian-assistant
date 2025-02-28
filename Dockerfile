FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/main .
COPY .env .

# Create a directory for your markdown files
RUN mkdir -p /app/data

EXPOSE 7777

CMD ["./main"]
