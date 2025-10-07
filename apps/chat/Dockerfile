# Base
FROM rust:1.90-slim-bookworm AS base
WORKDIR /app
RUN apt-get update && apt-get install -y libssl-dev pkg-config && rm -rf /var/lib/apt/lists/*
ENV SQLX_OFFLINE=true

# Builder
FROM base AS builder
WORKDIR /app
COPY Cargo.toml Cargo.lock ./
COPY src ./src
RUN cargo build --release

# Dev
FROM base AS dev
WORKDIR /app
COPY . .
RUN cargo install cargo-watch
CMD ["cargo", "watch", "-x", "run"]

# Prod
FROM debian:bookworm-slim AS prod
WORKDIR /app
RUN apt-get update && apt-get install -y libssl-dev ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/target/release/chat /usr/local/bin/chat
EXPOSE 8082
CMD ["chat"]
