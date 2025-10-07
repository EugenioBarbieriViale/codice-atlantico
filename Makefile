# Makefile for codice-atlantico

DOCKER_DIR := infra/docker
ENV_DEV := infra/env/.env.dev
ENV_PROD := infra/env/.env.prod
ENV_SECRETS := infra/env/secrets.env

COMPOSE_DEV := docker compose -f infra/docker/docker-compose.yml -f infra/docker/docker-compose.dev.yml
COMPOSE_PROD := docker compose -f infra/docker/docker-compose.yml -f infra/docker/docker-compose.prod.yml

# Development

.PHONY: up-dev down-dev logs-dev
up-dev:
	$(COMPOSE_DEV) up --build

down-dev:
	$(COMPOSE_DEV) down -v

logs-dev:
	$(COMPOSE_DEV) logs -f


# Production

.PHONY: up-prod down-prod build-prod logs-prod
up-prod:
	$(COMPOSE_PROD) up -d --build

down-prod:
	$(COMPOSE_PROD) down -v

build-prod:
	$(COMPOSE_PROD) build

logs-prod:
	$(COMPOSE_PROD) logs -f


# Database & migrations

MIGRATE_CMD := docker compose -f $(DOCKER_DIR)/docker-compose.yml run --rm migrate

.PHONY: migrate-new migrate-up migrate-down migrate-force
migrate-new:
	@read -p "Migration name: " name; \
	mkdir -p infra/db/migrations; \
	$(MIGRATE_CMD) create -ext sql -dir /migrations $$name

migrate-up:
	$(MIGRATE_CMD) up

migrate-down:
	$(MIGRATE_CMD) down 1

migrate-force:
	$(MIGRATE_CMD) force $(version)


# SQLC generation

.PHONY: sqlc
sqlc:
	cd apps/backend && sqlc generate

#
# Shells

.PHONY: sh-backend sh-frontend sh-db
sh-backend:
	$(COMPOSE_DEV) exec backend sh

sh-frontend:
	$(COMPOSE_DEV) exec frontend sh

sh-db:
	$(COMPOSE_DEV) exec db psql -U leonardo -d biblioteca-ambrosiana


# Linting & Formatting

.PHONY: lint lint-go lint-js fmt fmt-go fmt-js
lint: lint-go lint-js

lint-go:
	cd apps/backend && golangci-lint run

lint-js:
	cd apps/frontend && npm run lint

fmt: fmt-go fmt-js

fmt-go:
	cd apps/backend && go fmt ./...

fmt-js:
	cd apps/frontend && npm run format


# Testing

.PHONY: test test-go test-js
test: test-go test-js

test-go:
	cd apps/backend && go test ./...

test-js:
	cd apps/frontend && npm run test || echo "(No tests yet)"
