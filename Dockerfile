# Stage 1: Build Client (React)
FROM node:18 AS client-builder
WORKDIR /app
COPY ./Client/package*.json ./Client/
COPY ./Client/ ./
RUN npm install --prefix ./Client && npm run build --prefix ./Client

# Stage 2: Build Server (Go)
FROM golang:1.22.5 AS server-builder
WORKDIR /app
COPY ./Server/go.mod ./Server/go.sum ./
RUN go mod download
COPY ./Server/ ./
RUN go build -o /server

# Stage 3: Final Image
FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=server-builder /server /server
COPY --from=client-builder /app/Client/build /app/client
EXPOSE 8080
CMD ["/server"]
