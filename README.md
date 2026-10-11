# 🏨 Grand Oasis Hotel — AI Sales Agent & Fullstack Next.js Platform

Sistem pemesanan kamar hotel bintang lima modern yang mengintegrasikan **AI Sales Agent otonom berkinerja tinggi (Go & Gemini GenAI SDK)** dengan **Aplikasi Fullstack Web Portal & Receptionist Command Center (Next.js 15, React 19, Tailwind CSS, & NextAuth.js)**.

Proyek ini terinspirasi dari arsitektur berkinerja tinggi [`AI-Sales-Agent-Workshop-Go`](../AI-Sales-Agent-Workshop-Go), yang diadaptasi khusus untuk domain bisnis perhotelan dengan proteksi inventori *date-range capacity*, multi-channel support (WhatsApp & Web SSE), dan kontainerisasi multi-platform (**Docker & Podman**).

---

## 🌟 Fitur Unggulan Sistem

- **🧠 AI Sales Agent & Autonomous Loop (Rama Concierge)**:
  - Menggunakan SDK resmi [`google.golang.org/genai`](https://pkg.go.dev/google.golang.org/genai) dengan loop pemanggilan fungsi rekursif (*multi-turn agentic loop*).
  - Dilengkapi **6 Hotel Business Tools**: `search_rooms`, `check_room_availability`, `hold_booking`, `get_guest_bookings`, `cancel_booking`, dan `transfer_to_receptionist`.
  - Persona concierge ramah berstandar *hospitality* bintang lima dan pemformatan *bubble chat* alami menggunakan pemisah `[NEXT]`.
- **📅 Engine Inventori Kamar Berbasis Rentang Tanggal (Date-Range Capacity)**:
  - Kalkulasi ketersediaan unit kamar malam demi malam antara tanggal Check-In dan Check-Out.
  - Proteksi *race condition* dan *overbooking* simultan menggunakan **Mutex Registry per Room** dan transaksi atomik database.
- **⏳ Background Hold Expiry Worker**:
  - Worker otomatis (Ticker) berjalan setiap menit untuk membatalkan reservasi `HOLD` yang melewati batas waktu 2 jam tanpa verifikasi pembayaran, melepaskan kuota kamar kembali ke publik.
- **💬 Multi-Channel Ingestion (WhatsApp & Web Chat)**:
  - **WhatsApp**: Integrasi gateway Evolution API dengan **Debounce Buffer** (4 detik) agar chat beruntun disatukan sebelum diproses AI.
  - **Live Web Chat Widget**: Tertanam langsung di website tamu dengan streaming **Server-Sent Events (SSE)** berlatensi rendah.
- **🖥️ Fullstack Next.js 15 Portal & Command Center**:
  - **Landing Page Publik**: Hero section dengan pencarian ketersediaan kilat (*quick search*), katalog kartu kamar, dan fasilitas resor.
  - **Receptionist Command Center (`/admin`)**: Autentikasi NextAuth.js, live chat dua kolom, saklar **Takeover Mode (AI ON/OFF)**, form balasan manual staf, dan modal verifikasi bukti bayar.
- **🐳 Multi-Platform Containerization (Docker & Podman)**:
  - Otomasi skrip [`container.sh`](./container.sh) cerdas mendeteksi `podman`, `docker`, maupun compose tool di sistem Linux/macOS.

---

## 🏗️ Diagram Arsitektur Sistem

```
WhatsApp Tamu                                         Web Browser Tamu
      │                                                      │
      ▼ (Webhook)                                            ▼ (SSE Stream)
Evolution API Gateway                                Next.js App Router (Port 3000)
      │                                                      │
      ├───────────────────────┐              ┌───────────────┤
      │ POST /webhook         │              │ REST API / SSE│
      ▼                       ▼              ▼               ▼
┌────────────────────────────────────────────────────────────────────────────┐
│                             Go Backend Engine                              │
│                                (Port 8080)                                 │
│  - Gin HTTP Router & SSE Chunk Streamer                                    │
│  - WhatsApp Debounce Buffer & Bubble Splitter ([NEXT])                     │
│  - Gemini Agentic Loop (Google GenAI SDK)                                  │
│  - 6 Hotel Business Tools (GORM ORM)                                       │
│  - Atomic Inventory Concurrency Lock & Hold Expiration Ticker Worker       │
└─────────────────────────────────────┬──────────────────────────────────────┘
                                      │
                                      ▼
                        ┌───────────────────────────┐
                        │   PostgreSQL / SQLite     │
                        │ (Rooms, Bookings, Guests) │
                        └───────────────────────────┘
```

---

## 📚 Modul Dokumentasi Teknis Lengkap

Detail arsitektur, teknologi, dan alur data telah didokumentasikan ke dalam modul terpisah:

### 🔹 Modul Backend (Go Engine) — [`backend/docs/`](./backend/docs/)
1. [**Modul 01: Arsitektur Backend Go & Gin Server**](./backend/docs/01-arsitektur-backend.md)
2. [**Modul 02: Teknologi & Stack Dependensi Backend**](./backend/docs/02-teknologi-dan-stack.md)
3. [**Modul 03: Data Model & Skema Database (ERD)**](./backend/docs/03-data-model-dan-skema.md)
4. [**Modul 04: Engine Inventori Kamar & Proteksi Overbooking**](./backend/docs/04-engine-inventori-dan-concurrency.md)
5. [**Modul 05: Background Worker Pembersih Hold Expiry**](./backend/docs/05-hold-expiration-worker.md)
6. [**Modul 06: Gemini Agentic Loop & 6 Tools Bisnis Hotel**](./backend/docs/06-gemini-agentic-loop-dan-tools.md)
7. [**Modul 07: Flow Data & Multi-Channel Gateway (WhatsApp & Web SSE)**](./backend/docs/07-flow-data-dan-multi-channel.md)
8. [**Modul 08: Kontrak Admin REST API & CLI Simulator**](./backend/docs/08-admin-api-dan-cli-simulator.md)

### 🔹 Modul Frontend (Next.js 15 Fullstack) — [`frontend/docs/`](./frontend/docs/)
1. [**Modul 01: Arsitektur Next.js 15 App Router**](./frontend/docs/01-arsitektur-nextjs.md)
2. [**Modul 02: Teknologi, Styling, & Desain Resor Mewah**](./frontend/docs/02-teknologi-dan-styling.md)
3. [**Modul 03: Flow Data & Integrasi API Client**](./frontend/docs/03-flow-data-dan-api-client.md)
4. [**Modul 04: Live Web Chat Widget & SSE Streaming**](./frontend/docs/04-live-chat-widget-dan-sse.md)
5. [**Modul 05: Autentikasi NextAuth.js & Middleware Security**](./frontend/docs/05-autentikasi-dan-middleware.md)
6. [**Modul 06: Receptionist Command Center & Takeover UI**](./frontend/docs/06-admin-command-center.md)

---

## 🚀 Panduan Memulai Cepat

### Opsi 1: Menggunakan Script Container (`container.sh`)
Skrip helper `./container.sh` secara cerdas mendeteksi apakah komputer Anda menggunakan **Podman** atau **Docker**.

```bash
# 1. Pindah ke direktori project
cd AI-Booking-Hotel-System

# 2. Bangun image container Backend & Frontend
./container.sh build

# 3. Jalankan Database, Go Backend, dan Next.js Frontend
./container.sh up

# 4. Buka CLI simulator di dalam container
./container.sh cli

# 5. Matikan seluruh container saat selesai
./container.sh down
```

### Opsi 2: Menggunakan Docker Compose Mandiri (Backend / Frontend Saja)

Jika Anda ingin menjalankan backend atau frontend secara mandiri dalam container:

```bash
# Menjalankan Backend mandiri (Port 8080):
cd backend
cp .env.example .env   # Sesuaikan konfigurasi DB & API key
docker compose up -d

# Menjalankan Frontend mandiri (Port 3000):
cd frontend
cp .env.example .env   # Opsional: sesuaikan URL backend jika berbeda
docker compose up -d
```

### Opsi 3: Menjalankan Secara Lokal (Development)

**Persyaratan**: Go 1.25+ dan Node.js 22+/24+.

```bash
# Terminal 1: Backend Go (Port 8080)
cd AI-Booking-Hotel-System/backend
go run ./cmd/server

# Terminal 2: Frontend Next.js 15 (Port 3000)
cd AI-Booking-Hotel-System/frontend
npm run dev

# Terminal 3: Uji Coba Cepat Agen Rama di Terminal
cd AI-Booking-Hotel-System/backend
go run ./cmd/cli
```

---

## 🌐 Akses Aplikasi & Akun Default

* **Website Publik & Live Web Chat Widget**: [http://localhost:3000](http://localhost:3000)
* **Command Center Resepsionis & Admin**: [http://localhost:3000/admin](http://localhost:3000/admin)
  - Akun Admin (GM): `admin` / `admin123`
  - Akun Resepsionis: `receptionist` / `staff123`
* **Go REST API & Webhook**: [http://localhost:8080](http://localhost:8080)
