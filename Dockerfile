# Build Go backend
FROM golang:1.22.5-alpine AS server-build
WORKDIR /app
COPY Server/ . 
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main

# Build React frontend
FROM node:18-alpine AS client-build
WORKDIR /app
COPY Client/package*.json ./
RUN npm install
COPY Client/ . 
RUN npm run build

# Final stage
FROM alpine:latest
WORKDIR /app

# Install Node.js untuk `serve`
RUN apk add --no-cache nodejs npm
RUN npm install -g serve

# Copy backend binary dan frontend build
COPY --from=server-build /app/main /app/main
COPY --from=client-build /app/dist /app/dist

# Set environment variables
ENV PORT=8080

# Ekspos port untuk backend dan frontend
EXPOSE 8000 8080

# Jalankan backend di port 8080 dan frontend di port 8000
CMD ./main & serve -s /app/dist -l ${PORT:-8000}

