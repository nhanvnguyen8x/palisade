.PHONY: dev infra migrate api build

infra:
	docker compose up -d

migrate:
	docker compose exec -T postgres psql -U palisade -d palisade < scripts/migration/create-document.sql
	docker compose exec -T postgres psql -U palisade -d palisade < scripts/migration/phase2-knowledge-management.sql
	docker compose exec -T postgres psql -U palisade -d palisade < scripts/migration/phase3-knowledge-indexing.sql

api:
	cd apps/api && go run ./cmd/server

build:
	cd apps/api && go build -o bin/server ./cmd/server

dev: infra
	@echo "Waiting for services to start..."
	@sleep 3
	-$(MAKE) migrate
	$(MAKE) api

worker:
	@echo "Starting AI Runtime worker (apps/ai-runtime)..."
	cd apps/ai-runtime && python3 -m worker.main

worker-install:
	cd apps/ai-runtime && python3 -m pip install -r requirements.txt
