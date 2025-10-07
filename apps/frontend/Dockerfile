# Base
FROM node:20-alpine AS base
WORKDIR /app
RUN apk add --no-cache libc6-compat

# Dev
FROM base AS dev
COPY apps/frontend/package*.json ./
RUN npm ci
# Source will be mounted in dev
EXPOSE 3000
# Default command overridden by compose.dev

# Builder
FROM base AS builder
COPY apps/frontend/package*.json ./
RUN npm ci
COPY apps/frontend/ .
# Ensure SvelteKit uses adapter-node; build produces build/
RUN npm run build

# Prod
FROM node:20-alpine AS prod
WORKDIR /app
ENV NODE_ENV=production
# Only what's needed to run SSR
COPY --from=builder /app/build ./build
COPY --from=builder /app/package*.json ./
RUN npm ci --omit=dev
EXPOSE 3000
CMD ["node", "build/index.js"]
