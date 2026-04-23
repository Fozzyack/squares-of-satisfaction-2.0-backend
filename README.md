# Backend

Go API for authentication and habit management.

![Tiny Wins Landing](./tinywinslanding.png)

## Stack

- Go 1.26
- Chi router
- PostgreSQL
- Goose migrations (embedded and auto-run on startup)

## Prerequisites

- Go installed
- Docker (for local Postgres)

## Environment

Copy env values:

```bash
cp .env.example .env
```

Main variables:

- `ENVIRONMENT` (`dev` or `prod`)
- `DATABASE_URL`
- `COOKIE_NAME`
- `FRONTEND_URL` and `FRONTEND_URL_DEV`

## Run Locally

1) Start Postgres:

```bash
docker compose up -d
```

2) Start the API:

```bash
go run .
```

Server runs on `http://localhost:8800` by default.

Use a custom port:

```bash
go run . -port 8801
```

## Migrations

- SQL migrations live in `migrations/*.sql`.
- They are embedded in the binary and executed automatically at startup.

## Seed Data

Run the database seeder:

```bash
go run ./cmd/seed
```

This seeds:

- Clears existing DB data first
- 1 user (`john@example.com` / `password123`)
- 1 habit (`Drink Water`)
- Random `habit_entries` + `habit_daily_totals` for the last 365 days

## API Overview

Public routes:

- `GET /health`
- `POST /users` (create user + session cookie)
- `POST /users/login` (login + session cookie)

Session-protected routes:

- `GET /users` (current user)
- `POST /habits`
- `GET /habits`
- `GET /habits/{habitId}`
- `PUT /habits/{habitId}`

## Frontend Integration Notes

- CORS allows `http://localhost:3000` and `http://127.0.0.1:3000`.
- Session auth depends on a cookie named by `COOKIE_NAME`.
