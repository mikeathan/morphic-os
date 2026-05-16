# morphic-os

Morphic-OS is a Self-Synthesizing AI Operating System designed for high-performance local nodes. It dynamically adapts to tasks by forging, compiling, sandboxing, testing, and registering new Wasm microservices.

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Or, for local development: Go 1.24+, Node.js 18+

### Production Run (Docker)
To start the entire system (Go Backend, SQLite DB, Next.js Frontend) using the multi-stage Debian containers, simply run:
```bash
docker-compose up -d
```
The frontend will be available at `http://localhost:3000` and the backend at `http://localhost:8080`.
The database is persisted in the `./data` directory.

### Local IDE Debugging (Development)
During active development, you may want to run the application outside of Docker to attach a debugger (like Delve).

1. Set the environment variable:
   ```bash
   export MORPHIC_ENV=dev
   ```
2. Start the Backend:
   ```bash
   cd backend
   go run ./cmd/server
   ```
3. Start the Frontend:
   ```bash
   cd frontend
   npm run dev
   ```