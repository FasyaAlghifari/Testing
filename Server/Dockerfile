# Menggunakan image Golang untuk build server
FROM golang:1.22.5 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

# Image runtime untuk server
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
