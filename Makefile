SHELL := /usr/bin/env bash
COMPOSE := docker compose --env-file .env -f deployments/compose/compose.yaml

.PHONY: bootstrap doctor config up up-ai down restart status ps logs test clean

bootstrap:
	./scripts/bootstrap.sh

doctor:
	./scripts/doctor.sh

config:
	$(COMPOSE) config

up:
	$(COMPOSE) up -d --build

up-ai:
	COMPOSE_PROFILES=local-ai $(COMPOSE) -f deployments/compose/compose.local-ai.yaml up -d --build

down:
	$(COMPOSE) down

restart: down up

status:
	./scripts/status.sh

ps:
	$(COMPOSE) ps

logs:
	$(COMPOSE) logs -f --tail=200

test:
	cd apps/brain-api && go test ./...
	cd apps/knowledge-worker && python3 -m pytest -q || true

clean:
	$(COMPOSE) down -v --remove-orphans
