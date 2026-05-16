# ── Stage 1: build ──────────────────────────────────────────
FROM node:20-alpine AS builder

# Vite embeds env vars at build time via import.meta.env
ARG VITE_API_BASE_URL=https://api-node-995881892656.southamerica-east1.run.app
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

WORKDIR /app

COPY package.json ./
RUN npm install --frozen-lockfile
RUN chmod -R +x node_modules/.bin

COPY . .
RUN npm run build

# ── Stage 2: serve ──────────────────────────────────────────
FROM nginx:stable-alpine AS runner

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:8080/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
