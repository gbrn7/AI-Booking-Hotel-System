# Hotel AI Sales Agent & Next.js Fullstack Platform Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a complete, high-performance Hotel AI Sales Agent system in Go with multi-channel support (WhatsApp & Web Chat Widget), atomic date-range inventory locking, hold expiration worker, and a modern Next.js Fullstack Portal for guest booking, landing page, and receptionist takeover dashboard.

**Architecture:** Decoupled monorepo containing `backend/` (Go Gin, Google GenAI SDK, GORM PostgreSQL/SQLite fallback, WebSocket/SSE) and `frontend/` (Next.js 15 App Router, Tailwind CSS, Shadcn UI, NextAuth.js) with containerization support for both Docker and Podman via `container.sh` and `docker-compose.yml`.

**Tech Stack:**
- Backend: Go 1.25, Gin, GORM, PostgreSQL (`gorm.io/driver/postgres`), SQLite (`github.com/glebarez/sqlite` for quick local testing/mocking), `google.golang.org/genai`, Gorilla WebSocket.
- Frontend: Next.js 15 (App Router), TypeScript, Tailwind CSS, Lucide React, NextAuth.js (Auth.js), Server-Sent Events (SSE).
- DevOps & Runtime: Docker, Podman, `podman-compose` / `docker compose`, Bash automation (`container.sh`).

## Global Constraints
- Target workspace directory: `AI-Booking-Hotel-System/`.
- All Go packages must belong to module `github.com/raygbrn/ai-booking-hotel-system`.
- Preserved `ThoughtSignature` handling for Gemini models (`gemini-2.5-flash` / `gemini-1.5-flash`).
- Atomic room inventory calculation: ketersediaan dihitung per malam menginap dengan lock concurrency.
- Hold expiration background ticker runs every 1-2 minutes.
- Support both Docker and Podman seamlessly via `container.sh`.

---

## File Structure

```
AI-Booking-Hotel-System/
├── .env.example
├── container.sh
├── docker-compose.yml
├── backend/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── cmd/
│   │   ├── server/main.go
│   │   └── cli/main.go
│   ├── internal/
│   │   ├── agent/
│   │   │   ├── agent.go
│   │   │   ├── agent_test.go
│   │   │   ├── prompt.go
│   │   │   └── thought.go
│   │   ├── buffer/
│   │   │   ├── buffer.go
│   │   │   └── buffer_test.go
│   │   ├── config/
│   │   │   └── config.go
│   │   ├── db/
│   │   │   ├── db.go
│   │   │   └── seed.go
│   │   ├── handler/
│   │   │   ├── chat.go
│   │   │   ├── webhook.go
│   │   │   ├── admin.go
│   │   │   └── room.go
│   │   ├── memory/
│   │   │   ├── memory.go
│   │   │   └── takeover.go
│   │   ├── model/
│   │   │   └── models.go
│   │   ├── service/
│   │   │   ├── inventory.go
│   │   │   ├── inventory_test.go
│   │   │   ├── expiry.go
│   │   │   └── expiry_test.go
│   │   ├── tools/
│   │   │   ├── tools.go
│   │   │   └── tools_test.go
│   │   └── whatsapp/
│   │       ├── client.go
│   │       └── parser.go
│   └── prompts/
│       └── system.txt
└── frontend/
    ├── Dockerfile
    ├── package.json
    ├── tsconfig.json
    ├── next.config.ts
    ├── tailwind.config.ts
    └── src/
        ├── app/
        │   ├── layout.tsx
        │   ├── page.tsx
        │   ├── (auth)/login/page.tsx
        │   ├── admin/
        │   │   ├── layout.tsx
        │   │   ├── page.tsx
        │   │   ├── takeover/page.tsx
        │   │   ├── bookings/page.tsx
        │   │   └── rooms/page.tsx
        │   └── api/auth/[...nextauth]/route.ts
        ├── components/
        │   ├── chat/LiveChatWidget.tsx
        │   ├── admin/TakeoverChat.tsx
        │   ├── admin/BookingVerificationModal.tsx
        │   ├── landing/HeroSection.tsx
        │   ├── landing/RoomGrid.tsx
        │   └── ui/
        ├── lib/
        │   ├── api.ts
        │   └── auth.ts
        └── types/
            └── index.ts
```

---

### Task 1: Environment Configuration, Dockerfiles, and `container.sh`

**Files:**
- Create: `AI-Booking-Hotel-System/.env.example`
- Create: `AI-Booking-Hotel-System/container.sh`
- Create: `AI-Booking-Hotel-System/docker-compose.yml`
- Create: `AI-Booking-Hotel-System/backend/Dockerfile`
- Create: `AI-Booking-Hotel-System/frontend/Dockerfile`

**Interfaces:**
- Consumes: System bash, Docker, Podman.
- Produces: CLI script `./container.sh` supporting `build`, `up`, `up:all`, `down`, `logs`, `status`, `cli`.

- [ ] **Step 1: Write `.env.example`**

```env
# Application Environment
APP_ENV=development
PORT=8080
FRONTEND_PORT=3000
BACKEND_URL=http://localhost:8080
FRONTEND_URL=http://localhost:3000

# Google Gemini API
GEMINI_API_KEY=your_gemini_api_key_here
GEMINI_MODEL=gemini-2.5-flash
SYSTEM_PROMPT_PATH=./prompts/system.txt

# Database Configuration (PostgreSQL)
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=hotel_admin
DB_PASSWORD=hotel_secret_pass
DB_NAME=hotel_booking_db
DB_SSLMODE=disable

# Hold Booking Policy
HOLD_DURATION_MINUTES=120
EXPIRY_CHECK_INTERVAL_SECONDS=60

# WhatsApp / Evolution API
EVOLUTION_API_URL=http://localhost:8083
EVOLUTION_API_KEY=hotel_evolution_secret_key
EVOLUTION_INSTANCE_NAME=hotel-sales-agent
WA_DEBOUNCE_SECONDS=4

# NextAuth / JWT Secret
NEXTAUTH_SECRET=super_secret_hotel_nextauth_token_key_change_in_prod
NEXTAUTH_URL=http://localhost:3000
```

- [ ] **Step 2: Create `docker-compose.yml` with SELinux `:Z` flag and Podman compatibility**

```yaml
version: '3.8'

services:
  db:
    image: postgres:16-alpine
    container_name: hotel-db
    restart: unless-stopped
    environment:
      POSTGRES_USER: ${DB_USER:-hotel_admin}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-hotel_secret_pass}
      POSTGRES_DB: ${DB_NAME:-hotel_booking_db}
    ports:
      - "${DB_PORT:-5432}:5432"
    volumes:
      - ./data/postgres:/var/lib/postgresql/data:Z
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-hotel_admin} -d ${DB_NAME:-hotel_booking_db}"]
      interval: 5s
      timeout: 5s
      retries: 5

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: hotel-backend
    restart: unless-stopped
    ports:
      - "${PORT:-8080}:8080"
    environment:
      - APP_ENV=production
      - PORT=8080
      - DB_DRIVER=postgres
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=${DB_USER:-hotel_admin}
      - DB_PASSWORD=${DB_PASSWORD:-hotel_secret_pass}
      - DB_NAME=${DB_NAME:-hotel_booking_db}
      - DB_SSLMODE=disable
      - GEMINI_API_KEY=${GEMINI_API_KEY}
      - GEMINI_MODEL=${GEMINI_MODEL:-gemini-2.5-flash}
      - HOLD_DURATION_MINUTES=${HOLD_DURATION_MINUTES:-120}
      - EVOLUTION_API_URL=${EVOLUTION_API_URL:-http://evolution-api:8083}
      - EVOLUTION_API_KEY=${EVOLUTION_API_KEY}
      - EVOLUTION_INSTANCE_NAME=${EVOLUTION_INSTANCE_NAME:-hotel-sales-agent}
    depends_on:
      db:
        condition: service_healthy

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: hotel-frontend
    restart: unless-stopped
    ports:
      - "${FRONTEND_PORT:-3000}:3000"
    environment:
      - NEXT_PUBLIC_BACKEND_URL=http://localhost:8080
      - BACKEND_INTERNAL_URL=http://backend:8080
      - NEXTAUTH_SECRET=${NEXTAUTH_SECRET:-super_secret_hotel_nextauth_token_key_change_in_prod}
      - NEXTAUTH_URL=http://localhost:3000
    depends_on:
      - backend

  evolution-api:
    image: atendai/evolution-api:v2.1.2
    container_name: hotel-evolution-api
    restart: unless-stopped
    profiles:
      - full
    ports:
      - "8083:8083"
    environment:
      - SERVER_PORT=8083
      - AUTHENTICATION_API_KEY=${EVOLUTION_API_KEY:-hotel_evolution_secret_key}
    volumes:
      - ./data/evolution:/evolution/instances:Z
```

- [ ] **Step 3: Create `backend/Dockerfile` and `frontend/Dockerfile`**

`backend/Dockerfile`:
```dockerfile
# Stage 1: Build binary
FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates tzdata
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/cli ./cmd/cli

# Stage 2: Minimal runtime
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/bin/server /app/server
COPY --from=builder /app/bin/cli /app/cli
COPY --from=builder /app/prompts /app/prompts
EXPOSE 8080
ENTRYPOINT ["/app/server"]
```

`frontend/Dockerfile`:
```dockerfile
# Stage 1: Dependencies
FROM node:24-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm ci

# Stage 2: Builder
FROM node:24-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
ENV NEXT_TELEMETRY_DISABLED=1
RUN npm run build

# Stage 3: Runner
FROM node:24-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
EXPOSE 3000
CMD ["node", "server.js"]
```

- [ ] **Step 4: Create executable `container.sh` with dual Podman / Docker detection**

Implement `AI-Booking-Hotel-System/container.sh` with tool detection (`podman-compose`, `docker compose`, etc.), `ensure_data_dirs`, `cmd_up`, `cmd_down`, `cmd_build`, `cmd_cli`, and chmod +x.

- [ ] **Step 5: Verify `container.sh` runs `--help` and tests detection**

Run: `chmod +x container.sh && ./container.sh help`  
Expected output: Displays banner with detected engine and list of available commands.

- [ ] **Step 6: Commit**

```bash
git add .env.example docker-compose.yml container.sh backend/Dockerfile frontend/Dockerfile
git commit -m "chore: add docker, podman, and container.sh configuration"
```

---

### Task 2: Backend Database Models, GORM Setup, & Seeding

**Files:**
- Create: `AI-Booking-Hotel-System/backend/go.mod`
- Create: `AI-Booking-Hotel-System/backend/internal/model/models.go`
- Create: `AI-Booking-Hotel-System/backend/internal/db/db.go`
- Create: `AI-Booking-Hotel-System/backend/internal/db/seed.go`
- Test: `AI-Booking-Hotel-System/backend/internal/db/db_test.go`

**Interfaces:**
- Produces: `model.Room`, `model.Guest`, `model.Booking`, `model.Conversation`, `model.Message`, `model.User`, `db.InitDB(cfg)`, `db.SeedDefaultData(db)`.

- [ ] **Step 1: Write failing test in `backend/internal/db/db_test.go`**

```go
package db_test

import (
	"testing"
	"github.com/raygbrn/ai-booking-hotel-system/internal/db"
	"github.com/raygbrn/ai-booking-hotel-system/internal/model"
)

func TestInitSQLiteAndSeed(t *testing.T) {
	database, err := db.InitTestDB()
	if err != nil {
		t.Fatalf("failed to init test db: %v", err)
	}

	err = db.SeedDefaultData(database)
	if err != nil {
		t.Fatalf("failed to seed default data: %v", err)
	}

	var roomCount int64
	database.Model(&model.Room{}).Count(&roomCount)
	if roomCount == 0 {
		t.Errorf("expected seeded rooms, got 0")
	}

	var userCount int64
	database.Model(&model.User{}).Count(&userCount)
	if userCount == 0 {
		t.Errorf("expected seeded users, got 0")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd backend && go test ./internal/db/...`  
Expected: FAIL (packages/files do not exist yet).

- [ ] **Step 3: Implement `internal/model/models.go`, `internal/db/db.go`, and `internal/db/seed.go`**

Define GORM models matching Section 3 of design spec with JSON serialization, relations, and default seeded hotel rooms (Deluxe Ocean View, Grand Suite, Superior Room, Family Suite) with sample prices and total units, as well as an admin user (`admin` / `admin123`, bcrypt hashed).

- [ ] **Step 4: Run test to verify it passes**

Run: `cd backend && go test -v ./internal/db/...`  
Expected: PASS with 4 room types and default users seeded.

- [ ] **Step 5: Commit**

```bash
git add backend/
git commit -m "feat(backend): add gorm models, database initialization, and default seed data"
```

---

### Task 3: Hotel Inventory Engine & Concurrency Locking

**Files:**
- Create: `AI-Booking-Hotel-System/backend/internal/service/inventory.go`
- Test: `AI-Booking-Hotel-System/backend/internal/service/inventory_test.go`

**Interfaces:**
- Consumes: `model.Room`, `model.Booking`, GORM DB.
- Produces: `service.CheckAvailability(db, checkIn, checkOut, roomType string) ([]RoomAvailability, error)`
- Produces: `service.HoldBooking(db, req HoldBookingRequest) (*model.Booking, error)`

- [ ] **Step 1: Write failing test for Date-Range Capacity and Concurrent Overbooking**

```go
package service_test

import (
	"sync"
	"testing"
	"time"
	"github.com/raygbrn/ai-booking-hotel-system/internal/db"
	"github.com/raygbrn/ai-booking-hotel-system/internal/model"
	"github.com/raygbrn/ai-booking-hotel-system/internal/service"
)

func TestHoldBookingConcurrency(t *testing.T) {
	database, _ := db.InitTestDB()
	// Create a room with total_units = 1
	room := model.Room{
		Name: "Penthouse Suite",
		Type: "penthouse",
		BasePricePerNight: 2500000,
		CapacityGuests: 2,
		TotalUnits: 1,
		IsActive: true,
	}
	database.Create(&room)

	checkIn := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	checkOut := time.Now().Add(48 * time.Hour).Format("2006-01-02")

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(guestIdx int) {
			defer wg.Done()
			req := service.HoldBookingRequest{
				GuestName: "Guest",
				GuestPhone: "628111111",
				RoomID: room.ID,
				CheckInDate: checkIn,
				CheckOutDate: checkOut,
				TotalRooms: 1,
			}
			_, err := service.HoldBooking(database, req)
			if err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()

	if successCount != 1 {
		t.Fatalf("expected exactly 1 booking to succeed for 1 available unit, got %d", successCount)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd backend && go test -v ./internal/service -run TestHoldBookingConcurrency`  
Expected: FAIL.

- [ ] **Step 3: Implement `internal/service/inventory.go`**

Implement:
1. `CheckAvailability(db, checkIn, checkOut, roomType)`: Iterates every night between check-in and check-out, counts overlapping active bookings (`status IN ('HOLD', 'CONFIRMED', 'CHECKED_IN')`), returns remaining available units.
2. `HoldBooking(db, req)`: Uses DB transaction with row lock or room-level mutex, checks night-by-night availability, generates unique `booking_code` (e.g. `HTL-YYYYMM-XXXX`), sets `hold_expires_at = NOW() + 2 hours`, saves record with status `HOLD`.

- [ ] **Step 4: Run test to verify it passes**

Run: `cd backend && go test -v ./internal/service/...`  
Expected: PASS (atomic concurrency guarantees 1 booking).

- [ ] **Step 5: Commit**

```bash
git add backend/internal/service/
git commit -m "feat(backend): implement atomic date-range inventory engine and concurrency lock"
```

---

### Task 4: Automated Hold Expiration Background Worker

**Files:**
- Create: `AI-Booking-Hotel-System/backend/internal/service/expiry.go`
- Test: `AI-Booking-Hotel-System/backend/internal/service/expiry_test.go`

**Interfaces:**
- Produces: `service.StartHoldExpiryWorker(db *gorm.DB, interval time.Duration) (chan struct{}, error)`
- Produces: `service.ExpireStaleHoldBookings(db *gorm.DB) (int64, error)`

- [ ] **Step 1: Write failing test in `expiry_test.go`**

Test that a booking created with expired `hold_expires_at` (e.g. 1 minute ago) has its status changed to `CANCELLED` when `ExpireStaleHoldBookings` runs, and valid bookings remain `HOLD`.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd backend && go test -v ./internal/service -run TestExpireStaleHoldBookings`  
Expected: FAIL.

- [ ] **Step 3: Implement `internal/service/expiry.go`**

Implement `ExpireStaleHoldBookings` which executes:
```go
result := db.Model(&model.Booking{}).
    Where("status = ? AND hold_expires_at < ?", "HOLD", time.Now()).
    Update("status", "CANCELLED")
```
And `StartHoldExpiryWorker` with a `time.NewTicker(interval)` running in a background goroutine with graceful shutdown channel.

- [ ] **Step 4: Run test to verify it passes**

Run: `cd backend && go test -v ./internal/service -run TestExpireStaleHoldBookings`  
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add backend/internal/service/expiry*
git commit -m "feat(backend): add automated hold booking expiration worker"
```

---

### Task 5: 6 Hotel Business Tools & Gemini Agentic Loop

**Files:**
- Create: `AI-Booking-Hotel-System/backend/prompts/system.txt`
- Create: `AI-Booking-Hotel-System/backend/internal/tools/tools.go`
- Create: `AI-Booking-Hotel-System/backend/internal/agent/prompt.go`
- Create: `AI-Booking-Hotel-System/backend/internal/agent/thought.go`
- Create: `AI-Booking-Hotel-System/backend/internal/agent/agent.go`
- Test: `AI-Booking-Hotel-System/backend/internal/tools/tools_test.go`

**Interfaces:**
- Produces 6 tools: `search_rooms`, `check_room_availability`, `hold_booking`, `get_guest_bookings`, `cancel_booking`, `transfer_to_receptionist`.
- Produces: `agent.NewSalesAgent(apiKey, model, systemPrompt, db)`
- Produces: `agent.ProcessUserMessage(ctx, sessionID, userMsg) (string, error)`

- [ ] **Step 1: Write system prompt in `prompts/system.txt`**

Write "Rama" concierge persona: hospitality tone, upsell capabilities, tool call instructions, `[NEXT]` bubble split formatting for WhatsApp, and payment transfer instructions upon booking hold.

- [ ] **Step 2: Write failing test in `tools_test.go` for all 6 tools**

Test invocation of `SearchRooms`, `CheckRoomAvailability`, `HoldBooking`, `GetGuestBookings`, `CancelBooking`, and `TransferToReceptionist` against mock/test DB.

- [ ] **Step 3: Implement `internal/tools/tools.go`**

Implement tool definitions and handlers returning JSON responses formatted for `google.golang.org/genai` FunctionDeclaration / FunctionCall.

- [ ] **Step 4: Implement `internal/agent/agent.go` with Thought Signature preservation**

Adapt the autonomous tool calling loop from `AI-Sales-Agent-Workshop-Go`:
- Preserve `ThoughtSignature` so reasoning models like `gemini-2.5-flash` don't lose context.
- Support recursive tool resolution up to 5 hops.
- Return final synthesized assistant response.

- [ ] **Step 5: Run tests to verify tool handlers pass**

Run: `cd backend && go test -v ./internal/tools/...`  
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add backend/prompts/ backend/internal/tools/ backend/internal/agent/
git commit -m "feat(backend): implement 6 hotel business tools and gemini agentic loop"
```

---

### Task 6: Multi-Channel Gateway (WhatsApp Buffer & Web Chat SSE)

**Files:**
- Create: `AI-Booking-Hotel-System/backend/internal/buffer/buffer.go`
- Create: `AI-Booking-Hotel-System/backend/internal/memory/memory.go`
- Create: `AI-Booking-Hotel-System/backend/internal/memory/takeover.go`
- Create: `AI-Booking-Hotel-System/backend/internal/whatsapp/client.go`
- Create: `AI-Booking-Hotel-System/backend/internal/whatsapp/parser.go`
- Create: `AI-Booking-Hotel-System/backend/internal/handler/webhook.go`
- Create: `AI-Booking-Hotel-System/backend/internal/handler/chat.go`
- Test: `AI-Booking-Hotel-System/backend/internal/buffer/buffer_test.go`

**Interfaces:**
- Produces: `buffer.DebounceBuffer`, `memory.ConversationManager`, `handler.WhatsAppWebhookHandler`, `handler.WebChatSSEHandler`.

- [ ] **Step 1: Write failing test for debounce buffer**

Test that multiple rapid messages sent within debounce window (e.g. 50ms in test) are concatenated into a single combined message.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd backend && go test -v ./internal/buffer`  
Expected: FAIL.

- [ ] **Step 3: Implement `buffer.go`, `takeover.go`, and `webhook.go`**

- If `is_human_takeover == true`, incoming WhatsApp messages are persisted to DB without calling Gemini.
- WhatsApp replies parse `[NEXT]` into multiple message bubbles with simulated typing indicators.
- Web Chat endpoint `/api/chat/stream` supports Server-Sent Events (SSE) for streaming text and thought tokens to frontend widget.

- [ ] **Step 4: Run buffer tests and verify they pass**

Run: `cd backend && go test -v ./internal/buffer/...`  
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add backend/internal/buffer/ backend/internal/memory/ backend/internal/whatsapp/ backend/internal/handler/
git commit -m "feat(backend): implement multi-channel gateway, debounce buffer, and sse web chat"
```

---

### Task 7: REST API Handlers, CLI Simulator, & Server Entrypoint

**Files:**
- Create: `AI-Booking-Hotel-System/backend/internal/handler/admin.go`
- Create: `AI-Booking-Hotel-System/backend/internal/handler/room.go`
- Create: `AI-Booking-Hotel-System/backend/cmd/server/main.go`
- Create: `AI-Booking-Hotel-System/backend/cmd/cli/main.go`

**Interfaces:**
- Produces HTTP Endpoints:
  - `GET /api/rooms` - Public room catalog
  - `GET /api/availability` - Public availability check
  - `POST /api/chat/stream` - Guest Web Chat SSE
  - `POST /webhook` - WhatsApp Evolution API
  - `GET /api/admin/conversations` - List active chats (WA & Web)
  - `POST /api/admin/conversations/:id/takeover` - Toggle AI ON / OFF
  - `POST /api/admin/conversations/:id/reply` - Send manual staff message
  - `GET /api/admin/bookings` - List bookings with status filter
  - `POST /api/admin/bookings/:id/confirm` - Verify & confirm booking
- Produces CLI binary: `cmd/cli/main.go` for terminal testing.

- [ ] **Step 1: Implement `handler/admin.go` and `handler/room.go`**
- [ ] **Step 2: Implement `cmd/server/main.go` wiring all routes, GORM DB, expiry worker, and CORS middleware**
- [ ] **Step 3: Implement `cmd/cli/main.go` interactive console agent loop**
- [ ] **Step 4: Verify Go build succeeds**

Run: `cd backend && go build -o /dev/null ./cmd/server && go build -o /dev/null ./cmd/cli`  
Expected: Succeeded with exit code 0.

- [ ] **Step 5: Commit**

```bash
git add backend/cmd/ backend/internal/handler/
git commit -m "feat(backend): complete gin http server, admin apis, and cli simulator"
```

---

### Task 8: Frontend Project Initialization & Design System

**Files:**
- Create: `AI-Booking-Hotel-System/frontend/package.json`
- Create: `AI-Booking-Hotel-System/frontend/tsconfig.json`
- Create: `AI-Booking-Hotel-System/frontend/next.config.ts`
- Create: `AI-Booking-Hotel-System/frontend/tailwind.config.ts`
- Create: `AI-Booking-Hotel-System/frontend/src/app/globals.css`
- Create: `AI-Booking-Hotel-System/frontend/src/types/index.ts`
- Create: `AI-Booking-Hotel-System/frontend/src/lib/api.ts`

**Interfaces:**
- Produces: Next.js 15 App with Tailwind CSS, Lucide icons, `output: 'standalone'`, and type definitions for Rooms, Bookings, Conversations, Messages.

- [ ] **Step 1: Initialize Next.js dependencies and TypeScript configuration in `frontend/`**
- [ ] **Step 2: Setup Tailwind CSS & styling tokens (Luxury Hospitality theme: Slate, Gold/Amber accents)**
- [ ] **Step 3: Setup typed API client in `frontend/src/lib/api.ts` to call Go backend endpoints**
- [ ] **Step 4: Verify build compiles without errors**

Run: `cd frontend && npm run build`  
Expected: Standalone build successful.

- [ ] **Step 5: Commit**

```bash
git add frontend/
git commit -m "chore(frontend): initialize next.js 15 app router, tailwind css, and api client"
```

---

### Task 9: Frontend Public Landing Page & Room Showcase

**Files:**
- Create: `AI-Booking-Hotel-System/frontend/src/components/landing/HeroSection.tsx`
- Create: `AI-Booking-Hotel-System/frontend/src/components/landing/QuickSearch.tsx`
- Create: `AI-Booking-Hotel-System/frontend/src/components/landing/RoomCard.tsx`
- Create: `AI-Booking-Hotel-System/frontend/src/components/landing/RoomGrid.tsx`
- Create: `AI-Booking-Hotel-System/frontend/src/components/landing/Amenities.tsx`
- Create: `AI-Booking-Hotel-System/frontend/src/app/page.tsx`

**Interfaces:**
- Produces: Landing page with responsive hero search bar, room catalog cards with images and amenities tags, and booking trigger modals.

- [ ] **Step 1: Implement `QuickSearch.tsx` (Date-range picker check-in/out & guest count)**
- [ ] **Step 2: Implement `RoomCard.tsx` and `RoomGrid.tsx` with dynamic availability status**
- [ ] **Step 3: Compose `app/page.tsx` with Hero, RoomGrid, Amenities, and Hotel FAQ**
- [ ] **Step 4: Verify page renders and passes build check**

Run: `cd frontend && npm run build`  
Expected: Static/dynamic page build passes.

- [ ] **Step 5: Commit**

```bash
git add frontend/src/components/landing/ frontend/src/app/page.tsx
git commit -m "feat(frontend): build public hotel landing page and room showcase"
```

---

### Task 10: Frontend Floating Live Web Chat Widget

**Files:**
- Create: `AI-Booking-Hotel-System/frontend/src/components/chat/LiveChatWidget.tsx`
- Create: `AI-Booking-Hotel-System/frontend/src/components/chat/ChatMessage.tsx`
- Create: `AI-Booking-Hotel-System/frontend/src/components/chat/QuickChips.tsx`
- Modify: `AI-Booking-Hotel-System/frontend/src/app/layout.tsx`

**Interfaces:**
- Produces: Floating widget in bottom-right corner connecting to Go `/api/chat/stream` SSE, displaying typing indicators, Rama concierge avatar, and quick suggestion chips.

- [ ] **Step 1: Create `ChatMessage.tsx` supporting markdown formatting and booking status cards**
- [ ] **Step 2: Implement SSE stream handling in `LiveChatWidget.tsx` with session UUID storage in localStorage**
- [ ] **Step 3: Embed `LiveChatWidget` in `app/layout.tsx`**
- [ ] **Step 4: Verify widget functionality and component build**

Run: `cd frontend && npm run build`  
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add frontend/src/components/chat/ frontend/src/app/layout.tsx
git commit -m "feat(frontend): add live web chat widget with sse streaming and concierge persona"
```

---

### Task 11: Frontend Admin Authentication & Receptionist Takeover Center

**Files:**
- Create: `AI-Booking-Hotel-System/frontend/src/lib/auth.ts`
- Create: `AI-Booking-Hotel-System/frontend/src/app/api/auth/[...nextauth]/route.ts`
- Create: `AI-Booking-Hotel-System/frontend/src/app/(auth)/login/page.tsx`
- Create: `AI-Booking-Hotel-System/frontend/src/components/admin/TakeoverChat.tsx`
- Create: `AI-Booking-Hotel-System/frontend/src/app/admin/takeover/page.tsx`
- Create: `AI-Booking-Hotel-System/frontend/src/middleware.ts`

**Interfaces:**
- Produces: Protected `/admin` routes via NextAuth.js JWT session, real-time two-column live chat inbox (WhatsApp & Web), AI Takeover Switch (ON/OFF), and staff manual reply box.

- [ ] **Step 1: Configure NextAuth Credentials Provider in `src/lib/auth.ts` validating against Go admin auth endpoint**
- [ ] **Step 2: Implement Login UI in `app/(auth)/login/page.tsx`**
- [ ] **Step 3: Implement `TakeoverChat.tsx` with conversation list, AI toggle button, and manual message sender**
- [ ] **Step 4: Create Next.js route protection middleware for `/admin/*`**
- [ ] **Step 5: Verify build compiles without errors**

Run: `cd frontend && npm run build`  
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add frontend/src/lib/auth.ts frontend/src/app/api/auth/ frontend/src/app/(auth)/ frontend/src/components/admin/TakeoverChat.tsx frontend/src/app/admin/takeover/ frontend/src/middleware.ts
git commit -m "feat(frontend): implement admin authentication and live takeover center"
```

---

### Task 12: Frontend Reservation Management & Room Inventory CMS

**Files:**
- Create: `AI-Booking-Hotel-System/frontend/src/components/admin/BookingVerificationModal.tsx`
- Create: `AI-Booking-Hotel-System/frontend/src/app/admin/bookings/page.tsx`
- Create: `AI-Booking-Hotel-System/frontend/src/app/admin/rooms/page.tsx`
- Create: `AI-Booking-Hotel-System/frontend/src/app/admin/page.tsx`

**Interfaces:**
- Produces:
  - Admin Overview KPI metrics (Occupancy rate, Pending Hold bookings, Today's check-ins).
  - Booking table with filter badges (`HOLD`, `CONFIRMED`, `CANCELLED`).
  - Verification dialog: Preview transfer proof & button "Verifikasi & Konfirmasi Booking".
  - Room management table: Edit price, toggle active status, adjust units.

- [ ] **Step 1: Implement `BookingVerificationModal.tsx` triggering Go `/api/admin/bookings/:id/confirm`**
- [ ] **Step 2: Build `app/admin/bookings/page.tsx` with filter tabs and search by guest/code**
- [ ] **Step 3: Build `app/admin/rooms/page.tsx` for room price and inventory adjustment**
- [ ] **Step 4: Build `app/admin/page.tsx` (Dashboard summary cards)**
- [ ] **Step 5: Verify build passes**

Run: `cd frontend && npm run build`  
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add frontend/src/components/admin/ frontend/src/app/admin/
git commit -m "feat(frontend): implement reservation verification center and room cms"
```

---

### Task 13: End-to-End Multi-Platform Container Verification

**Files:**
- Test scripts: `AI-Booking-Hotel-System/container.sh`

**Verification Actions:**
- [ ] **Step 1: Test container build using Podman / Docker**

Run: `./container.sh build`  
Expected: Both `hotel-backend` and `hotel-frontend` images build successfully without errors.

- [ ] **Step 2: Test starting services using `container.sh up`**

Run: `./container.sh up`  
Expected: Postgres, Backend, and Frontend containers start and pass healthchecks.

- [ ] **Step 3: Test CLI Simulator inside container**

Run: `./container.sh cli` (or run CLI locally)  
Expected: Terminal simulator opens, greets user as Rama Concierge, checks availability when asked, and holds a booking.

- [ ] **Step 4: Test Hold Expiration Worker**

Verify that test hold bookings past expiration time are automatically transitioned to `CANCELLED`.

- [ ] **Step 5: Final Git commit & tag**

```bash
git add .
git commit -m "chore: complete end-to-end hotel booking platform and container verification"
```
