# Stage 1: Build Client (React)
FROM node:18 AS client-builder
WORKDIR /app
COPY ./Client/package*.json ./
RUN npm install
COPY ./Client/ ./
RUN npm run build

# Stage 2: Build Server (Go)
FROM golang:1.22.5 AS server-builder
WORKDIR /app
COPY ./Server/go.mod ./Server/go.sum ./
RUN go mod download
COPY ./Server/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o server

# Stage 3: Final Image
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=server-builder /app/server ./
COPY --from=client-builder /app/dist ./client
EXPOSE 8080
CMD ["./server"]
