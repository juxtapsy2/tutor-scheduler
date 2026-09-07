.PHONY: up down build logs migrate seed test dev

up:
	docker compose up --build -d

down:
	docker compose down

build:
	docker compose build

logs:
	docker compose logs -f

migrate:
	docker exec -i tutor-scheduler-pg psql -U brightpath -d brightpath < backend/migrations/001_schema.sql
	docker exec -i tutor-scheduler-pg psql -U brightpath -d brightpath < backend/migrations/002_seed.sql

test:
	cd backend && go test ./... -count=1

dev:
	docker compose up -d postgres
	@echo "Waiting for postgres..."
	@sleep 3
	cd backend && go run ./cmd/server/ &
	cd frontend && npm run dev
