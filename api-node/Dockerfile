FROM node:20-alpine AS deps

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci --omit=dev

FROM node:20-alpine AS runner

RUN addgroup -g 10001 appgroup && adduser -D -u 10001 -G appgroup app

WORKDIR /app

RUN apk add --no-cache tini curl

ENV NODE_ENV=production
ENV APP_HOST=0.0.0.0
ENV APP_PORT=3001

COPY --from=deps --chown=app:appgroup /app/node_modules ./node_modules
COPY --chown=app:appgroup package.json package.json
COPY --chown=app:appgroup src ./src

USER app

EXPOSE 3001

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -fsS http://127.0.0.1:3001/health

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["node", "src/index.js"]
