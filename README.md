# SmartSpa Admin

Step 1 delivers the backend skeleton with SQLite models and migrations.

## Project Structure
- server/ — Go backend (Gin + GORM + SQLite)
  - main.go — entrypoint with Gin server and health route
  - go.mod — module definition and dependencies
  - internal/models — core GORM entities (Member, Technician, ServiceItem, Appointment, Schedule, FissionLog)
  - internal/db — SQLite connection and automigration
  - internal/response — unified API envelope helpers
- web/ — (reserved) Vue 3 + Vite frontend scaffold (to be implemented in later steps)
- spa_management.db — SQLite data file (created on first run)
- AGENTS.md — prompt and implementation guide

## Run (backend)
1. `cd server`
2. `go mod tidy`
3. `go run main.go`

The server listens on :8080 with a health probe at `/health`.

## Run (frontend)
1. Install Bun (>=1.1.0).
2. `cd web`
3. `bun install`
4. `bun run dev` (default at :5173)

Environment: set `VITE_API_BASE_URL` to point to the Gin backend (defaults to `http://localhost:8080`).
