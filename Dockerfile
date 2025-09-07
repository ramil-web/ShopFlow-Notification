FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o notification ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/notification .
COPY .env .
EXPOSE 8082
CMD ["./notification"]
