# Menggunakan image untuk backend Golang
FROM golang:1.22.5 AS server-build

# Set working directory khusus untuk backend
WORKDIR /app/Server

# Copy hanya folder Server dan build backend
COPY Server/ /app/Server
RUN go mod download
RUN go build -o main .

# Menggunakan image untuk frontend Node.js
FROM node:18 AS client-build

# Set working directory khusus untuk frontend
WORKDIR /app/Client

# Copy hanya folder Client dan install dependencies
COPY Client/ /app/Client
RUN npm install
RUN npm run build

# Menggunakan image untuk menjalankan aplikasi
FROM golang:1.20

# Copy binary backend dan hasil build frontend
COPY --from=server-build /app/Server/main /app/main
COPY --from=client-build /app/Client/build /app/Client/build

# Set working directory untuk aplikasi
WORKDIR /app

# Ekspos port 8080 (atau port backend yang Anda gunakan)
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
