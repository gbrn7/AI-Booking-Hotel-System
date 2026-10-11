# Standalone Docker Compose for Backend and Frontend Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Provide standalone docker-compose setups and template environments inside both `backend/` and `frontend/` directories without relying on the root compose file.

**Architecture:** Independent `docker-compose.yml` files in `backend/` and `frontend/` using bridge networking with `host.docker.internal` host gateway resolution for connecting to external databases and backend services running on the host or other networks.

**Tech Stack:** Docker Compose v3.8, Go 1.23, Next.js 14/15, Alpine Linux.

## Global Constraints
- Do not modify or break root `docker-compose.yml`.
- Standalone compose files must use version `'3.8'`.
- Maintain existing Dockerfile configurations and standard environment variable names.
- Provide clean `.env.example` templates in both `backend/` and `frontend/`.

---

### Task 1: Backend Standalone Docker Compose & Env Template

**Files:**
- Create: `backend/.env.example`
- Create: `backend/docker-compose.yml`

**Interfaces:**
- Consumes: `backend/Dockerfile`, existing Go backend env conventions (`PORT`, `DB_*`, `GEMINI_*`, `EVOLUTION_*`).
- Produces: Runnable standalone backend container on port 8080 connecting to host/external database.

- [ ] **Step 1: Create backend/.env.example**

Write `backend/.env.example`:
```env
# Application Environment
APP_ENV=development
PORT=8080

# Google Gemini API
GEMINI_API_KEY=your_gemini_api_key_here
GEMINI_MODEL=gemini-2.5-flash
SYSTEM_PROMPT_PATH=./prompts/system.txt

# Database Configuration (PostgreSQL / SQLite)
# When running backend in Docker and PostgreSQL runs on your host machine, set DB_HOST=host.docker.internal
DB_DRIVER=postgres
DB_HOST=host.docker.internal
DB_PORT=5432
DB_USER=hotel_admin
DB_PASSWORD=hotel_secret_pass
DB_NAME=hotel_booking_db
DB_SSLMODE=disable

# Hold Booking Policy
HOLD_DURATION_MINUTES=120
EXPIRY_CHECK_INTERVAL_SECONDS=60

# WhatsApp / Evolution API
EVOLUTION_API_URL=http://host.docker.internal:8083
EVOLUTION_API_KEY=hotel_evolution_secret_key
EVOLUTION_INSTANCE=hotel-sales-agent
EVOLUTION_INSTANCE_NAME=hotel-sales-agent
WA_DEBOUNCE_SECONDS=4
```

- [ ] **Step 2: Create backend/docker-compose.yml**

Write `backend/docker-compose.yml`:
```yaml
version: '3.8'

services:
  backend:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: hotel-backend
    restart: unless-stopped
    ports:
      - "${PORT:-8080}:8080"
    extra_hosts:
      - "host.docker.internal:host-gateway"
    env_file:
      - .env
    environment:
      - APP_ENV=${APP_ENV:-production}
      - PORT=${PORT:-8080}
      - DB_DRIVER=${DB_DRIVER:-postgres}
      - DB_HOST=${DB_HOST:-host.docker.internal}
      - DB_PORT=${DB_PORT:-5432}
      - DB_USER=${DB_USER:-hotel_admin}
      - DB_PASSWORD=${DB_PASSWORD:-hotel_secret_pass}
      - DB_NAME=${DB_NAME:-hotel_booking_db}
      - DB_SSLMODE=${DB_SSLMODE:-disable}
      - GEMINI_API_KEY=${GEMINI_API_KEY}
      - GEMINI_MODEL=${GEMINI_MODEL:-gemini-2.5-flash}
      - HOLD_DURATION_MINUTES=${HOLD_DURATION_MINUTES:-120}
      - EVOLUTION_API_URL=${EVOLUTION_API_URL:-http://host.docker.internal:8083}
      - EVOLUTION_API_KEY=${EVOLUTION_API_KEY}
      - EVOLUTION_INSTANCE=${EVOLUTION_INSTANCE:-${EVOLUTION_INSTANCE_NAME:-hotel-sales-agent}}
      - EVOLUTION_INSTANCE_NAME=${EVOLUTION_INSTANCE:-${EVOLUTION_INSTANCE_NAME:-hotel-sales-agent}}
```

- [ ] **Step 3: Verify backend compose configuration syntax**

Run: `docker-compose -f backend/docker-compose.yml config`
Expected: Valid compose config with service `backend` outputted without errors.

- [ ] **Step 4: Commit changes**

```bash
git add backend/.env.example backend/docker-compose.yml
git commit -m "feat(backend): add standalone docker-compose and env template"
```

---

### Task 2: Frontend Standalone Docker Compose & Env Template

**Files:**
- Create: `frontend/.env.example`
- Create: `frontend/docker-compose.yml`

**Interfaces:**
- Consumes: `frontend/Dockerfile`, Next.js environment variable conventions (`NEXT_PUBLIC_BACKEND_URL`, `BACKEND_INTERNAL_URL`, `NEXTAUTH_*`).
- Produces: Runnable standalone frontend container on port 3000 connecting to backend.

- [ ] **Step 1: Create frontend/.env.example**

Write `frontend/.env.example`:
```env
# Frontend Port
FRONTEND_PORT=3000

# Backend API URLs
# NEXT_PUBLIC_BACKEND_URL is used by the browser client
NEXT_PUBLIC_BACKEND_URL=http://localhost:8080
# BACKEND_INTERNAL_URL is used by Next.js Server-Side Rendering (SSR) & API routes inside the container
BACKEND_INTERNAL_URL=http://host.docker.internal:8080

# NextAuth Configuration
NEXTAUTH_SECRET=super_secret_hotel_nextauth_token_key_change_in_prod
NEXTAUTH_URL=http://localhost:3000
```

- [ ] **Step 2: Create frontend/docker-compose.yml**

Write `frontend/docker-compose.yml`:
```yaml
version: '3.8'

services:
  frontend:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: hotel-frontend
    restart: unless-stopped
    ports:
      - "${FRONTEND_PORT:-3000}:3000"
    extra_hosts:
      - "host.docker.internal:host-gateway"
    env_file:
      - .env
    environment:
      - NEXT_PUBLIC_BACKEND_URL=${NEXT_PUBLIC_BACKEND_URL:-http://localhost:8080}
      - BACKEND_INTERNAL_URL=${BACKEND_INTERNAL_URL:-http://host.docker.internal:8080}
      - NEXTAUTH_SECRET=${NEXTAUTH_SECRET:-super_secret_hotel_nextauth_token_key_change_in_prod}
      - NEXTAUTH_URL=${NEXTAUTH_URL:-http://localhost:3000}
```

- [ ] **Step 3: Verify frontend compose configuration syntax**

Run: `docker-compose -f frontend/docker-compose.yml config`
Expected: Valid compose config with service `frontend` outputted without errors.

- [ ] **Step 4: Commit changes**

```bash
git add frontend/.env.example frontend/docker-compose.yml
git commit -m "feat(frontend): add standalone docker-compose and env template"
```

---

### Task 3: Integration Verification & Documentation

**Files:**
- Modify: `README.md`

**Interfaces:**
- Consumes: `backend/docker-compose.yml`, `frontend/docker-compose.yml`, root `docker-compose.yml`.
- Produces: Documented and verified standalone compose commands for developers.

- [ ] **Step 1: Update README.md with standalone Docker Compose instructions**

Add clear guidance in `README.md` under the Docker section explaining how to run backend or frontend independently.

- [ ] **Step 2: Verify all compose files**

Run:
```bash
docker-compose -f docker-compose.yml config > /dev/null && echo "Root compose OK"
docker-compose -f backend/docker-compose.yml config > /dev/null && echo "Backend compose OK"
docker-compose -f frontend/docker-compose.yml config > /dev/null && echo "Frontend compose OK"
```
Expected: All three print "OK".

- [ ] **Step 3: Commit documentation**

```bash
git add README.md
git commit -m "docs: document standalone docker compose usage for backend and frontend"
```
