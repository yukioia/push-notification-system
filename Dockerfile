# Stage 1: build
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o push-server ./src/cmd/server

# Stage 2: final image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/push-server .
COPY ./web ./web

EXPOSE 8080

CMD ["./push-server"]
