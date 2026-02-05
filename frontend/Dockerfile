# ===== BUILD STAGE =====
FROM node:20-alpine AS builder

WORKDIR /app

# Kopiraj package fajlove
COPY package*.json ./

RUN npm install

# Kopiraj ceo frontend
COPY . .

# Build produkcione fajlove
RUN npm run build

# ===== RUNTIME STAGE =====
FROM nginx:alpine

# Obriši default nginx config
RUN rm /etc/nginx/conf.d/default.conf

# Dodaj custom nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Kopiraj build fajlove
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
