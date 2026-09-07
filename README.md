# Bright Path Scheduler

Scheduling system for Bright Path Learning Centre in Da Nang. Prevents conflicting bookings for students, tutors, and rooms.

## Tech Stack

- **Backend:** Go 1.25, Gin, pgx, PostgreSQL 16
- **Frontend:** React 19, TypeScript, Vite, Tailwind CSS
- **Infrastructure:** Docker Compose

## Prerequisites

- Docker & Docker Compose
- Go 1.25+
- Node.js 22+

## Quick Start (Docker)

```bash
# Start all services (postgres + backend + frontend) [Easiest]
make up

# Run database migrations and seed data
make migrate
```

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080/api

## Local Development

```bash
# Start only postgres in Docker, run backend & frontend on host
make dev
```

- Frontend: http://localhost:5173 (proxies /api to :8080)
- Backend: http://localhost:8080

## Commands

| Command | Description |
|---------|-------------|
| `make up` | Build and start all services |
| `make down` | Stop all services |
| `make build` | Rebuild Docker images |
| `make migrate` | Run schema + seed SQL |
| `make test` | Run Go tests |
| `make dev` | Local dev (postgres in Docker, backend/frontend on host) |
| `make logs` | Tail service logs |

## API

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/appointments` | Create a booking |
| `GET` | `/api/appointments` | List appointments (optional: `?date=`, `?studentId=`, `?tutorId=`, `?roomId=`) |
| `GET` | `/api/students` | List students |
| `GET` | `/api/tutors` | List tutors |
| `GET` | `/api/rooms` | List rooms |

## Project Structure

```
├── backend/
│   ├── cmd/server/          # Entry point
│   ├── internal/
│   │   ├── handler/         # HTTP handlers
│   │   ├── application/     # Business logic
│   │   ├── domain/          # Types and errors
│   │   ├── repository/      # Database queries
│   │   └── db/              # Connection setup
│   └── migrations/          # SQL schema + seed
├── frontend/
│   └── src/
│       ├── features/
│       │   ├── booking/         # BookingForm
│       │   └── appointments/    # AppointmentList
│       └── shared/              # API client, types
├── docker-compose.yml
├── .env
└── Makefile
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `POSTGRES_USER` | `brightpath` | Database user |
| `POSTGRES_PASSWORD` | `brightpath` | Database password |
| `POSTGRES_DB` | `brightpath` | Database name |
| `POSTGRES_HOST` | `localhost` | Database host |
| `POSTGRES_PORT` | `5432` | Database port |
| `POSTGRES_SSLMODE` | `disable` | SSL mode |
