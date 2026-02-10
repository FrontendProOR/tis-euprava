FROM node:22-alpine

WORKDIR /app

# instalacija dependecy-ja
COPY package*.json ./
RUN npm ci

# kopiraj ostatak
COPY . .

EXPOSE 5173

# Vite mora da sluša na 0.0.0.0 u dockeru
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "5173"]
