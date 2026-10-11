# Design Spec: Standalone Docker Compose Setup for Backend and Frontend

- **Date:** 2026-10-11
- **Status:** Approved
- **Scope:** Provide standalone `docker-compose.yml` and `.env.example` configurations within both `backend/` and `frontend/` directories.

---

## 1. Context and Motivation

The project currently has a root-level `docker-compose.yml` that orchestrates the entire stack (Database, Backend, Frontend, and Evolution API). For modular local development, testing, and independent deployment, developers need the ability to run either the backend or frontend services standalone directly from their respective directories (`backend/` and `frontend/`).

---

## 2. Architecture and Technical Design

### 2.1 Backend Standalone Docker Compose (`backend/docker-compose.yml`)

- **Service Name:** `backend` (container name: `hotel-backend`)
- **Build Context:** `.` (using `backend/Dockerfile`)
- **Port Mapping:** `${PORT:-8080}:8080`
- **Host Resolution:**
  ```yaml
  extra_hosts:
    - "host.docker.internal:host-gateway"
  ```
  This allows backend running inside the container to connect to an external PostgreSQL database running on the host machine or another container (via `DB_HOST=host.docker.internal` or host IP).
- **Environment Management:**
  - Loads `./.env` if present (`env_file: [.env]`).
  - Sets safe defaults for variables such as `PORT=8080`, `DB_DRIVER=${DB_DRIVER:-postgres}`, `DB_HOST=${DB_HOST:-host.docker.internal}`, `DB_PORT=${DB_PORT:-5432}`, `GEMINI_API_KEY=${GEMINI_API_KEY}`, etc.
- **Support Files:**
  - `backend/.env.example`: Provides a clean template for backend developers running standalone.

### 2.2 Frontend Standalone Docker Compose (`frontend/docker-compose.yml`)

- **Service Name:** `frontend` (container name: `hotel-frontend`)
- **Build Context:** `.` (using `frontend/Dockerfile`)
- **Port Mapping:** `${FRONTEND_PORT:-3000}:3000`
- **Host Resolution:**
  ```yaml
  extra_hosts:
    - "host.docker.internal:host-gateway"
  ```
  Allows Next.js server-side rendering (SSR) and API routes inside the container to communicate with the Go backend running on host via `http://host.docker.internal:8080`.
- **Environment Management:**
  - Browser-side: `NEXT_PUBLIC_BACKEND_URL=${NEXT_PUBLIC_BACKEND_URL:-http://localhost:8080}`
  - Server-side: `BACKEND_INTERNAL_URL=${BACKEND_INTERNAL_URL:-http://host.docker.internal:8080}`
  - NextAuth: `NEXTAUTH_URL=${NEXTAUTH_URL:-http://localhost:3000}`, `NEXTAUTH_SECRET=${NEXTAUTH_SECRET:-super_secret_hotel_nextauth_token_key_change_in_prod}`
  - Loads `./.env` or `./.env.local` if present.
- **Support Files:**
  - `frontend/.env.example`: Configures frontend environment variables.

---

## 3. Directory Structure Changes

```
AI-Booking-Hotel-System/
├── backend/
│   ├── Dockerfile
│   ├── docker-compose.yml      # NEW: Standalone backend compose
│   ├── .env.example            # NEW: Backend template env
│   └── ...
├── frontend/
│   ├── Dockerfile
│   ├── docker-compose.yml      # NEW: Standalone frontend compose
│   ├── .env.example            # NEW: Frontend template env
│   └── ...
└── docs/superpowers/specs/
    └── 2026-10-11-docker-compose-backend-frontend-design.md
```

---

## 4. Verification & Testing Strategy

1. **Syntax & Config Validation:**
   - Run `docker compose config` inside `backend/` to verify YAML validity and variable substitution.
   - Run `docker compose config` inside `frontend/` to verify YAML validity and variable substitution.
2. **Build Validation:**
   - Verify `docker compose build` succeeds in both directories without breaking existing Dockerfile workflows.
3. **Integration Check:**
   - Verify root `docker-compose.yml` remains completely intact and unaffected.
