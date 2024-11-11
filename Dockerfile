# Build Go backend
FROM golang:1.20-alpine AS backend-builder
WORKDIR /app
COPY Server/ .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main

# Build React frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app
COPY Client/package*.json ./
RUN npm install
COPY Client/ .
RUN npm run build

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=backend-builder /app/main .
COPY --from=frontend-builder /app/dist ./dist
COPY Server/.env .

EXPOSE 8080
CMD ["./main"]