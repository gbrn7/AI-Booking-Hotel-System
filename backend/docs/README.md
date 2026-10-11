# 📚 Modul Dokumentasi Teknis Backend (Go Engine)

Selamat datang di dokumentasi teknis mendalam untuk **Go Backend & AI Sales Agent Engine** pada sistem **Grand Oasis Hotel Booking Platform**.

Backend ini dirancang dengan prinsip konkurensi tinggi, integritas data atomik, dan arsitektur *agentic loop* berbasis **Google GenAI SDK** dengan multi-channel support (WhatsApp & Web Chat SSE).

---

## 🗺️ Daftar Modul Pembelajaran & Referensi Teknis

| Modul | Topik | Deskripsi Utama |
|---|---|---|
| [**Modul 01**](./01-arsitektur-backend.md) | **Arsitektur Sistem & Struktur Layer** | Pola arsitektur *Clean / Layered Architecture*, pembagian tanggung jawab (`cmd/`, `internal/`), dan lifecycle HTTP request. |
| [**Modul 02**](./02-teknologi-dan-stack.md) | **Teknologi & Dependensi** | Ulasan mendalam Go 1.25+, Gin Framework, GORM ORM, PostgreSQL & SQLite fallback, serta Google GenAI SDK. |
| [**Modul 03**](./03-data-model-dan-skema.md) | **Data Model & Relasi Database** | Skema tabel PostgreSQL, relasi GORM (`Room`, `Guest`, `Booking`, `Conversation`, `Message`, `User`), dan State Machine booking. |
| [**Modul 04**](./04-engine-inventori-dan-concurrency.md) | **Engine Inventori & Concurrency Lock** | Algoritma kalkulasi ketersediaan kamar per tanggal menginap, pencegahan overbooking, dan muteks sinkronisasi thread. |
| [**Modul 05**](./05-hold-expiration-worker.md) | **Background Hold Expiry Worker** | Siklus hidup background worker (Ticker) untuk auto-cancel reservasi berstatus `HOLD` yang melewati batas waktu 2 jam. |
| [**Modul 06**](./06-gemini-agentic-loop-dan-tools.md) | **AI Agentic Loop & 6 Hotel Tools** | Arsitektur pemanggilan fungsi autonomous Gemini, preservasi *thought process*, dan 6 tools bisnis perhotelan. |
| [**Modul 07**](./07-flow-data-dan-multi-channel.md) | **Flow Data & Multi-Channel Gateway** | Alur data WhatsApp webhook (Debounce Buffer + Bubble Split) dan Web Chat SSE (*Server-Sent Events*), serta Human Takeover switch. |
| [**Modul 08**](./08-admin-api-dan-cli-simulator.md) | **Admin REST API & CLI Simulator** | Spesifikasi kontrak endpoint HTTP Admin, autentikasi JWT, dan panduan penggunaan simulator terminal interaktif (`cmd/cli`). |

---

## 📂 Struktur Direktori Backend

```
backend/
├── cmd/
│   ├── server/               # Entrypoint HTTP Server Gin & API Router
│   └── cli/                  # Simulator AI Concierge Interaktif di Terminal
├── internal/
│   ├── agent/                # Gemini GenAI SDK Loop, Prompt Compiler, & Reasoning
│   ├── buffer/               # Debounce Buffer pesan chat WhatsApp
│   ├── config/               # Pemuatan Environment Variables (.env)
│   ├── db/                   # Koneksi GORM PostgreSQL / SQLite & Seeding
│   ├── handler/              # HTTP Controller (Room, Chat SSE, Webhook, Admin)
│   ├── memory/               # Conversation Session Tracker & Takeover Manager
│   ├── model/                # GORM Struct Models & Relasi Database
│   ├── service/              # Inventory Concurrency Engine & Hold Expiry Worker
│   ├── tools/                # 6 Hotel Business Tools (GORM query execution)
│   └── whatsapp/             # Klien Evolution API & Parser Webhook Payload
├── prompts/
│   └── system.txt            # System Prompt Persona AI Rama (Concierge Specialist)
├── docs/                     # Dokumentasi teknis & modul pembelajaran
├── Dockerfile                # Multi-stage Dockerfile Go (Binary statis alpine)
├── go.mod / go.sum           # Manajemen dependensi Go
└── hotel.db                  # Database SQLite lokal (jika mode dev/offline)
```
