# Hotel AI Sales Agent & Next.js Fullstack Platform - Design Specification

**Date:** 2026-10-11  
**Status:** Approved  
**Target Path:** `AI-Booking-Hotel-System/`  
**Reference:** `AI-Sales-Agent-Workshop-Go`

---

## 1. Executive Summary & Goals

Aplikasi ini adalah platform pemesanan kamar hotel modern yang menggabungkan:
1. **AI Sales Agent & Reservation Engine (Go Backend)**: Mengadaptasi arsitektur kinerja tinggi dari `AI-Sales-Agent-Workshop-Go`, mengotomatisasi konsultasi kamar, rekomendasi, pengecekan ketersediaan tanggal, dan pembuatan reservasi berstatus *HOLD*.
2. **Multi-Channel Communication**: Mendukung interaksi tamu melalui WhatsApp (Evolution API) dan Live Web Chat Widget terintegrasi.
3. **Fullstack Next.js Portal (Frontend & Dashboard)**:
   - *Public Landing Page*: Menampilkan profil hotel, katalog tipe kamar, fasilitas, dan quick availability search bar.
   - *Live Web Chat Widget*: Widget percakapan real-time untuk tamu berinteraksi langsung dengan AI Concierge.
   - *Admin & Receptionist Dashboard*: Autentikasi berbasis NextAuth.js, live chat takeover mode, verifikasi pembayaran manual, dan pengelolaan ketersediaan kamar.

---

## 2. System Architecture

Arsitektur sistem menggunakan pola **Decoupled Monorepo** di dalam folder `AI-Booking-Hotel-System/`:

```
                  ┌──────────────────────────────────────────────┐
                  │                 Tamu Hotel                   │
                  └──────────────┬────────────────┬──────────────┘
                                 │                │
                        (WhatsApp Chat)     (Web Browser)
                                 │                │
                                 ▼                ▼
                     ┌────────────────────┐  ┌─────────────────────────┐
                     │   Evolution API    │  │   Next.js App Router    │
                     │  (WhatsApp Gateway)│  │ (Landing Page, Widget,  │
                     └─────────┬──────────┘  │    Admin Dashboard)     │
                               │             └────────────┬────────────┘
                        (HTTP Webhook)                    │ (REST / SSE API)
                               ▼                          ▼
                     ┌─────────────────────────────────────────────────┐
                     │               Go Backend Engine                 │
                     │  - Gin HTTP Router & SSE/WS Web Chat Stream     │
                     │  - Gemini Agentic Loop (Google GenAI SDK)       │
                     │  - 6 Hotel Business Tools                       │
                     │  - Debounce Buffer & Takeover Manager           │
                     │  - Room Inventory Concurrency & Hold Expiry     │
                     └────────────────────────┬────────────────────────┘
                                              │ (GORM ORM)
                                              ▼
                                 ┌────────────────────────┐
                                 │  PostgreSQL Database   │
                                 └────────────────────────┘
```

### Directory Structure
```
AI-Booking-Hotel-System/
├── backend/                  # Go Gin Backend & AI Agent Engine
│   ├── cmd/
│   │   ├── server/           # Entrypoint HTTP Server & Webhook Gateway
│   │   └── cli/              # Simulator CLI Interaktif untuk pengujian lokal
│   ├── internal/
│   │   ├── agent/            # Gemini SDK loop, system prompt, thought signature handler
│   │   ├── buffer/           # WhatsApp debounce buffer per nomor telepon
│   │   ├── config/           # Pemuatan environment variables
│   │   ├── db/               # Inisialisasi PostgreSQL GORM & auto-migration / seed
│   │   ├── handler/          # HTTP & SSE handlers (Chat, Webhook, Admin, Rooms, Bookings)
│   │   ├── memory/           # Riwayat chat conversation & takeover state tracker
│   │   ├── model/            # GORM Structs (Room, Guest, Booking, Conversation, Message, User)
│   │   ├── service/          # Inventory concurrency engine & Hold expiration background worker
│   │   ├── tools/            # 6 Hotel Business Tools
│   │   └── whatsapp/         # Client Evolution API
│   ├── prompts/              # System prompt & Persona AI Concierge ("Rama")
│   ├── go.mod / go.sum
│   └── Dockerfile
├── frontend/                 # Next.js Fullstack (App Router)
│   ├── src/
│   │   ├── app/
│   │   │   ├── (public)/     # Landing page, Room details, Booking success info
│   │   │   ├── (auth)/login/ # Admin / Receptionist login page
│   │   │   ├── admin/        # Dashboard overview, Live Takeover chat, Booking table, Rooms
│   │   │   └── api/auth/     # NextAuth.js route handler
│   │   ├── components/
│   │   │   ├── chat/         # Floating Live Web Chat Widget & SSE client
│   │   │   ├── admin/        # Takeover chat UI, Booking verification modals
│   │   │   └── ui/           # Shadcn UI primitives
│   │   ├── lib/              # API client, NextAuth configuration, helper utils
│   │   └── types/            # TypeScript interfaces
│   ├── package.json
│   ├── tailwind.config.ts
│   └── Dockerfile            # Multi-stage Dockerfile Next.js (Standalone output)
├── docs/
│   └── superpowers/
│       └── specs/            # Design specifications
├── container.sh              # Skrip otomasi multi-platform (Docker & Podman detection)
├── docker-compose.yml        # PostgreSQL, Go Backend, Next.js Frontend, Evolution API
└── .env.example
```

---

## 3. Data Model (PostgreSQL)

### 3.1 `rooms`
* `id` (UUID / Serial, PK)
* `name` (VARCHAR, e.g. "Deluxe Ocean View")
* `type` (VARCHAR, `standard` | `deluxe` | `suite` | `family_suite`)
* `base_price_per_night` (DECIMAL / INT)
* `capacity_guests` (INT)
* `total_units` (INT, jumlah total unit fisik untuk tipe kamar ini)
* `description` (TEXT)
* `amenities` (JSONB, e.g. `["King Bed", "Ocean View", "Bathtub", "Free WiFi", "Breakfast Included"]`)
* `image_urls` (JSONB)
* `is_active` (BOOLEAN, default: true)
* `created_at`, `updated_at` (TIMESTAMP)

### 3.2 `guests`
* `id` (UUID / Serial, PK)
* `name` (VARCHAR)
* `phone_number` (VARCHAR, unique identifier untuk WhatsApp & Web Chat)
* `email` (VARCHAR, nullable)
* `created_at`, `updated_at` (TIMESTAMP)

### 3.3 `bookings`
* `id` (UUID / Serial, PK)
* `booking_code` (VARCHAR, UNIQUE, e.g. "HTL-202610-8X9Y")
* `guest_id` (FK -> `guests.id`)
* `room_id` (FK -> `rooms.id`)
* `check_in_date` (DATE)
* `check_out_date` (DATE)
* `total_nights` (INT)
* `total_rooms` (INT, default: 1)
* `total_price` (DECIMAL / INT)
* `status` (VARCHAR: `HOLD`, `CONFIRMED`, `CANCELLED`, `CHECKED_IN`, `CHECKED_OUT`)
* `hold_expires_at` (TIMESTAMP, default: `created_at + 2 hours`)
* `payment_proof_url` (VARCHAR, nullable)
* `special_requests` (TEXT, nullable)
* `confirmed_by_user_id` (FK -> `users.id`, nullable)
* `confirmed_at` (TIMESTAMP, nullable)
* `created_at`, `updated_at` (TIMESTAMP)

### 3.4 `conversations` & `messages`
* `conversations`:
  * `id` (VARCHAR, PK, session key: e.g. `wa:628123456789` atau `web:<uuid>`)
  * `guest_id` (FK -> `guests.id`, nullable)
  * `channel` (VARCHAR: `whatsapp` | `web`)
  * `is_human_takeover` (BOOLEAN, default: false)
  * `last_message_at` (TIMESTAMP)
* `messages`:
  * `id` (UUID / Serial, PK)
  * `conversation_id` (FK -> `conversations.id`)
  * `sender` (VARCHAR: `user` | `assistant` | `agent_tool` | `human_agent`)
  * `content` (TEXT)
  * `thought_process` (TEXT, nullable)
  * `created_at` (TIMESTAMP)

### 3.5 `users`
* `id` (UUID / Serial, PK)
* `username` (VARCHAR, UNIQUE)
* `email` (VARCHAR, UNIQUE)
* `password_hash` (VARCHAR)
* `name` (VARCHAR)
* `role` (VARCHAR: `admin` | `receptionist`)
* `created_at`, `updated_at` (TIMESTAMP)

---

## 4. AI Sales Agent Specifications

### 4.1 Persona & Sales Flow
* **Nama Agent**: "Rama" — Hotel Concierge & Sales Specialist untuk "Grand Oasis Hotel".
* **Karakter**: Ramah, sopan (*hospitality tone*), sigap memberikan rekomendasi kamar sesuai kebutuhan (bisnis, keluarga, liburan romantis), proaktif menawarkan opsi sarapan/layanan antar-jemput.
* **Flow Penjualan**:
  1. *Eksplorasi*: Memahami tanggal menginap, jumlah tamu, dan preferensi tipe kamar.
  2. *Cek Ketersediaan*: Memanggil tool `check_room_availability`.
  3. *Presentasi Opsi*: Menampilkan kamar yang cocok beserta keunggulan fasilitas dan harga.
  4. *Hold Booking*: Meminta konfirmasi nama & nomor kontak, lalu memanggil tool `hold_booking`.
  5. *Petunjuk Pembayaran*: Menyajikan kode booking, nominal transfer, rekening resmi bank hotel, dan menginfokan batas waktu hold 2 jam.

### 4.2 Hotel Business Tools (Go GORM)
1. **`search_rooms`**:
   - Filter kamar berdasarkan rentang harga, tipe (`standard`, `deluxe`, `suite`), dan kapasitas tamu.
2. **`check_room_availability`**:
   - Parameter: `check_in_date` (YYYY-MM-DD), `check_out_date` (YYYY-MM-DD), `room_type` (opsional).
   - Menghitung ketersediaan unit kamar setiap malam dalam rentang tanggal.
3. **`hold_booking`**:
   - Parameter: `guest_name`, `guest_phone`, `room_id`, `check_in_date`, `check_out_date`, `total_rooms`, `special_requests`.
   - Melakukan pengecekan kapasitas atomik dengan database lock.
   - Membuat record `bookings` dengan status `HOLD` dan `hold_expires_at = NOW() + 2 hours`.
4. **`get_guest_bookings`**:
   - Parameter: `phone_number` atau `booking_code`.
   - Mengambil status booking aktif dan riwayat pemesanan tamu.
5. **`cancel_booking`**:
   - Parameter: `booking_code`, `phone_number`, `reason`.
   - Membatalkan booking jika status masih `HOLD` (atau sesuai ketentuan) dan mengembalikan unit ke inventori.
6. **`transfer_to_receptionist`**:
   - Parameter: `reason`.
   - Mengubah status percakapan menjadi `is_human_takeover = true` dan memicu alert real-time ke admin dashboard.

---

## 5. Inventory Engine & Concurrency Control

### 5.1 Perhitungan Ketersediaan Berbasis Rentang Tanggal
Kamar tersedia pada rentang $[D_{in}, D_{out})$ jika dan hanya jika untuk setiap malam $D \in [D_{in}, D_{out}-1]$:
$$\text{available\_units}(D) = \text{total\_units} - \text{count\_booked\_units}(D) \ge \text{requested\_units}$$
di mana $\text{count\_booked\_units}(D)$ menghitung seluruh booking pada kamar tersebut dengan status $\in \{\text{'HOLD'}, \text{'CONFIRMED'}, \text{'CHECKED\_IN'}\}$ yang melingkupi malam $D$.

### 5.2 Concurrency Protection
* Pada Go service `HoldBooking`, operasi dijalankan dalam transaksi PostgreSQL dengan `SELECT ... FOR UPDATE` pada baris ketersediaan atau menggunakan in-memory mutex terkoordinasi per `room_id`.
* Mencegah overbooking jika 2 request reservasi untuk unit terakhir masuk secara bersamaan dari WhatsApp dan Web Chat.

### 5.3 Automated Hold Expiration Worker
* Background worker berjalan setiap 2 menit:
  ```sql
  UPDATE bookings
  SET status = 'CANCELLED'
  WHERE status = 'HOLD' AND hold_expires_at < NOW();
  ```
* Unit kamar otomatis terlepas dan siap dipesan kembali.

---

## 6. Frontend Next.js & Dashboard Specifications

### 6.1 Landing Page & Guest Experience
* **Hero Search Banner**: Input tanggal Check-in, Check-out, jumlah tamu dengan tombol "Cek Ketersediaan".
* **Room Showcase**: Grid kartu kamar dengan foto, badge fasilitas, kapasitas, harga per malam, dan tombol detail.
* **Floating Live Web Chat Widget**:
  - Posisi: Bottom-right screen dengan bubble floating.
  - State: Minimize / Expanded.
  - Fitur: Quick suggestion chips (*"Cek ketersediaan weekend ini"*, *"Tipe kamar terbaik untuk keluarga"*), streaming balasan agen, dan indikator pengetikan.

### 6.2 Admin & Receptionist Dashboard (`/admin`)
* **Security & Auth**: NextAuth.js Credentials Provider dengan JWT session, proteksi route via Next.js Middleware.
* **Live Chat & Takeover Center**:
  - Kolom kiri: List percakapan aktif (WhatsApp & Web) dengan badge status AI ON / TAKEOVER.
  - Kolom kanan: Live transcript percakapan, tombol saklar **Takeover (AI ON / OFF)**, dan form balasan manual.
* **Booking & Verification Center**:
  - Tabel reservasi dengan filter status (`HOLD`, `CONFIRMED`, `CANCELLED`).
  - Modal verifikasi pembayaran: Memeriksa bukti transfer, mencocokkan nominal, dan menekan tombol **"Verifikasi & Konfirmasi"** yang mengubah status menjadi `CONFIRMED`.
* **Room Management**:
  - CRUD tipe kamar, penyesuaian harga harian, dan pengaturan kuota total unit.

---

## 7. CLI Simulator (`cmd/cli`)
* Menyediakan simulator CLI interaktif berbasis Go terminal (seperti di `AI-Sales-Agent-Workshop-Go`).
* Memungkinkan developer menguji persona agen, pemanggilan tool, dan alur booking langsung via keyboard tanpa dependensi WhatsApp atau browser.

---

## 8. Multi-Platform Containerization (Docker & Podman)

### 8.1 Dockerfiles
* **`backend/Dockerfile`**:
  - Multi-stage build Go dengan binary statis (`CGO_ENABLED=0 GOOS=linux`).
  - Runtime minimal `alpine:3.20` atau `distroless/static-debian12` dengan CA-certificates & tzdata Asia/Jakarta.
* **`frontend/Dockerfile`**:
  - Multi-stage build Node.js (Alpine).
  - Menggunakan fitur Next.js standalone (`output: "standalone"` di `next.config.mjs`) untuk ukuran image ringan (~100MB).

### 8.2 Multi-Service `docker-compose.yml`
* **`db`**: PostgreSQL 16 Alpine, port `5432`, volume persisten `./data/postgres:/var/lib/postgresql/data:Z` (dukungan SELinux Podman flag `:Z`).
* **`backend`**: Go Gin Server & Agent Engine, port `8080`, depend on `db`.
* **`frontend`**: Next.js App, port `3000`, depend on `backend`.
* **`evolution-api` (opsional profile)**: Evolution API v2 untuk gateway WhatsApp lokal, port `8083`.

### 8.3 Skrip Otomasi `container.sh`
Mengadaptasi `container.sh` dari `AI-Sales-Agent-Workshop-Go`:
* Mendeteksi otomatis tool container: `podman-compose`, `podman compose`, `docker compose`, atau `docker-compose`.
* Mendeteksi engine: `podman` atau `docker`.
* Menangani volume permission rootless Podman (label `:Z` dan `userns: keep-id`).
* Menyediakan perintah praktis:
  - `./container.sh build` : Build image backend & frontend
  - `./container.sh up` : Jalankan seluruh stack (DB + Backend + Frontend)
  - `./container.sh up:full` : Jalankan stack lengkap termasuk Evolution API
  - `./container.sh down` : Matikan seluruh container
  - `./container.sh logs` : Streaming log container
  - `./container.sh cli` : Jalankan interactive CLI simulator di dalam container backend
  - `./container.sh status` : Cek healthcheck & container status

---

## 9. Verification & Success Criteria
1. Agent mampu merekomendasikan kamar sesuai kriteria pengguna via chat.
2. Tool calling (`check_room_availability`, `hold_booking`) tereksekusi dengan benar tanpa halusinasi inventori.
3. Overbooking dicegah saat kuota unit kamar habis pada rentang tanggal yang diminta.
4. Auto-cancelation berjalan ketika batas waktu status `HOLD` terlampaui.
5. Tamu dapat chat via WhatsApp atau Web Chat Widget.
6. Resepsionis dapat melakukan takeover chat dan mengonfirmasi pembayaran reservasi melalui Next.js Admin Dashboard.
7. Aplikasi dapat di-build dan dijalankan dengan mulus baik di **Docker** maupun **Podman** melalui `./container.sh up`.

