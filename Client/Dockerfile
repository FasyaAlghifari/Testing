# Menggunakan Node untuk membangun aplikasi React
FROM node:16 AS build
WORKDIR /app
COPY . .
RUN npm install
RUN npm run build

# Menggunakan Nginx untuk menyajikan aplikasi React
FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
