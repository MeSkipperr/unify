# Stage 1: build
FROM node:25-slim AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

RUN npm run build

# Stage 2: production
FROM node:25-slim

WORKDIR /app

COPY --from=builder /app ./

ENV NODE_ENV=production

EXPOSE 3000

CMD ["npm", "run", "start"]