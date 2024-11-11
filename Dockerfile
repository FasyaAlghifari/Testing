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
FROM node:18-alpine
WORKDIR /app

# Copy binary backend dan hasil build frontend
COPY --from=server-build /app/main ./
COPY --from=client-build /app/dist ./dist
COPY --from=client-build /app/package*.json ./
COPY --from=client-build /app/node_modules ./node_modules

# Set environment variables
ENV PORT=8080
ENV NODE_ENV=production

EXPOSE 8080

# Buat script untuk menjalankan kedua service
COPY <<EOF /app/start.sh
#!/bin/sh
npm start &
./main
EOF

RUN chmod +x /app/start.sh
CMD ["/app/start.sh"]