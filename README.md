# PotTrack

PotTrack is a hobby plant tracking application (per pot) from seeding to harvest.

This repository contains a **Go/Gin** backend and a **SvelteKit + TailwindCSS** frontend. The database is PostgreSQL using Indonesian table names and column names as specified in the PRD.

## Repository Structure

```
backend/            # Go server, migrations, tests
frontend/           # SvelteKit application
PRD_POT_TRACK.txt   # Product requirements
README.md           # This file
```

## Prerequisites

- Go 1.20+
- Node.js 18+
- Docker & docker-compose (for Postgres)

## Database

The development database is `pot_track_dev`. A docker-compose file is provided below.

```yaml
# docker-compose.yml
version: '3.8'
services:
  db:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: pot_track_dev
    ports:
      - "5432:5432"
    volumes:
      - ./backend/data:/var/lib/postgresql/data
```
```

Start the database:

```bash
cd backend
docker-compose up -d
```

The Go server will read `DATABASE_URL` from the environment (e.g. `postgres://postgres:password@localhost:5432/pot_track_dev?sslmode=disable`).

Migrations are run automatically on startup using `github.com/golang-migrate/migrate` and the SQL files located in `backend/migrations`.

## Backend

1. Navigate to `backend/`.
2. Run `go mod tidy` to download dependencies.
3. Build or run:
   ```bash
   go run ./cmd/app
   ```
4. The API listens on port 8080 by default. Use the endpoints defined in the PRD (e.g. `POST /api/auth/register`).

To run unit tests:

```bash
go test ./...
```

## Frontend

1. Navigate to `frontend/`.
2. Install dependencies:
   ```bash
   npm install
   ```
3. Run dev server:
   ```bash
   npm run dev -- --open
   ```
4. The app will open at `http://localhost:5173` (or the port shown).
   Create a `.env` file in the `frontend` folder if you need to specify the backend base URL:
   ```env
   VITE_API_BASE=http://localhost:8080
   ```
   The `apiFetch` helper will prefix requests with this value.

The UI is mobile-friendly and uses Indonesian navigation labels as per the requirements.

## Development Notes

- Backend uses JWT authentication, bcrypt password hashing, and simple input validation with Gin's binding tags.
- Database schema uses Indonesian table/column names and corresponds to the PRD.
- Frontend currently has placeholders for each page; you can extend fetch logic to interact with the API.

## Further Work

- Implement full business logic (stage computation, scheduling, reminders, photo upload to object storage, etc.).
- Add email notification support.
- Expand frontend with forms and real data binding.

Happy gardening! 🌱
