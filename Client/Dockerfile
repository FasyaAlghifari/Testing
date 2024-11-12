# Stage 1: Build Client (React)
FROM node:18 AS client-builder
WORKDIR /app
COPY ./Client/package*.json ./Client/
RUN npm install --prefix ./Client
COPY ./Client/ ./Client
RUN npm run build --prefix ./Client
