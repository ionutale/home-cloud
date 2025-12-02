FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o home-cloud main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/home-cloud .
COPY --from=builder /app/public ./public

# Create directories for uploads and thumbnails
RUN mkdir -p uploads thumbnails

EXPOSE 8080 2121

CMD ["./home-cloud"]
