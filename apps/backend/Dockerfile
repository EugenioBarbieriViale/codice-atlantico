# Base
FROM golang:1.25-alpine AS base
WORKDIR /app
RUN apk add --no-cache build-base git

# Dev
FROM base AS dev
# Hot reload tool (Air)
RUN go install github.com/air-verse/air@latest
# Pre-copy mod files for better caching
COPY apps/backend/go.mod apps/backend/go.sum ./
RUN go mod download
# Source comes from volume in dev; still copy to allow 'air' config presence
COPY apps/backend/.air.toml ./.air.toml
# Default command is set in compose.dev
EXPOSE 8080

# Builder
FROM base AS builder
COPY apps/backend/go.mod apps/backend/go.sum ./
RUN go mod download
COPY apps/backend/ .
# Adjust main path if different
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server ./cmd/server

# Prod
FROM gcr.io/distroless/base-debian12:nonroot AS prod
WORKDIR /srv
COPY --from=builder /out/server /srv/server
ENV PORT=8080
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/srv/server"]
