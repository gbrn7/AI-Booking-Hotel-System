# Modul 01: Arsitektur Backend Go & Gin Server

## 1. Filosofi & Pola Arsitektur

Backend sistem **Grand Oasis Hotel** dibangun menggunakan bahasa **Go (Golang)** dengan pendekatan **Layered Clean Architecture**. Arsitektur ini memisahkan secara ketat tanggung jawab antara protokol input (HTTP / Webhook / CLI), logika bisnis inti (Domain Services & Tools), state persisten (Database Models), dan orkestrasi kecerdasan buatan (Gemini Agentic Loop).

### Diagram Komponen Backend
```
                       ┌───────────────────────────────┐
                       │        External Inputs        │
                       │ (Web Client, WhatsApp, CLI)   │
                       └───────────────┬───────────────┘
                                       │
                                       ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          HTTP Transport Layer                               │
│  - Gin Engine Router (CORS, JSON Binding)                                   │
│  - Handlers: RoomHandler, ChatHandler (SSE), WebhookHandler, AdminHandler   │
└───────────────┬───────────────────────────────┬─────────────────────────────┘
                │                               │
                ▼                               ▼
┌──────────────────────────────┐ ┌─────────────────────────────────────────────┐
│      Multi-Channel Ingestion │ │           Autonomous AI Layer               │
│ - Debounce Buffer            │ │ - Agentic Loop (google.golang.org/genai)    │
│ - Bubble Splitter ([NEXT])   │ │ - Thought Signature Preserver              │
│ - Takeover State Gate        │ │ - Prompt Renderer (prompts/system.txt)      │
└───────────────┬──────────────┘ └──────────────────────┬──────────────────────┘
                │                                       │
                └───────────────────┬───────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            Domain Service Layer                             │
│  - Hotel Inventory Service: CheckAvailability, HoldBooking (Atomic Lock)    │
│  - Background Expiry Worker: Periodical Auto-Cancellation                   │
│  - 6 Hotel Business Tools: Search, Availability, Hold, Get, Cancel, Handover│
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Persistence Layer (GORM)                            │
│  - PostgreSQL 16 (Production) / SQLite WAL (Dev/Test)                       │
│  - Models: Room, Guest, Booking, Conversation, Message, User                │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Struktur Direktori & Tanggung Jawab

| Paket | Path | Tanggung Jawab |
|---|---|---|
| **Entrypoints** | `cmd/server/main.go`<br>`cmd/cli/main.go` | Inisialisasi dependensi, pemuatan config, registrasi handler ke router Gin, dan penyiapan CLI terminal simulator. |
| **Configuration** | `internal/config/config.go` | Membaca variable lingkungan (`.env` dan os environment) serta menyediakan *fallback defaults*. |
| **Database** | `internal/db/db.go`<br>`internal/db/seed.go` | Menginisialisasi koneksi GORM (PostgreSQL atau SQLite in-memory untuk testing), migrasi skema tabel, dan seeding data default. |
| **Domain Models** | `internal/model/models.go` | Definisi *struct* entitas bisnis, relasi foreign key, dan konstanta status. |
| **Domain Services** | `internal/service/inventory.go`<br>`internal/service/expiry.go` | Logika bisnis kapasitas kamar berbasis tanggal, locking konkurensi, dan worker pembersih status HOLD kedaluwarsa. |
| **AI Agent** | `internal/agent/agent.go`<br>`internal/agent/prompt.go` | Loop rekursif pemanggilan tool Gemini, evaluasi konteks tanggal, dan preservasi penalaran AI. |
| **AI Tools** | `internal/tools/tools.go` | Definisi skema JSON Schema Function Declaration dan eksekusi SQL/GORM untuk setiap tool. |
| **Channels & Buffer** | `internal/buffer/buffer.go`<br>`internal/whatsapp/` | Penanganan debounce chat masuk berurutan dari WhatsApp dan komunikasi balik via Evolution API. |
| **Conversation Memory** | `internal/memory/memory.go`<br>`internal/memory/takeover.go` | Penyimpanan sliding window histori chat dan pengatur saklar Takeover staf hotel. |
| **HTTP Handlers** | `internal/handler/` | Pengontrol endpoint HTTP REST, Server-Sent Events (SSE) streaming chat, dan Webhook WhatsApp. |

---

## 3. Siklus Hidup Request (Request Lifecycle)

### A. Alur Tamu Web (Live Chat SSE)
1. Tamu mengirim pesan dari browser ke `POST /api/chat/stream`.
2. `ChatHandler` memvalidasi payload, memastikan session ID, dan menyimpan pesan user ke tabel `messages`.
3. Handler memeriksa `TakeoverManager`:
   - Jika `is_human_takeover == true`, handler langsung mengirim respons notifikasi bahwa obrolan ditangani staf resepsionis tanpa memanggil Gemini AI.
   - Jika takeover tidak aktif, handler mengambil 10 histori percakapan terakhir dan meneruskannya ke `agent.RunAgent`.
4. Agent loop berinteraksi dengan model Gemini. Bila model meminta function call (misal: `check_room_availability`), tool dieksekusi secara lokal terhadap database, hasilnya dikembalikan ke model hingga model menghasilkan balasan final.
5. Balasan dipecah berdasarkan tag `[NEXT]` dan distreaming langsung ke browser menggunakan format SSE (`data: {"type": "bubble", ...}`).

### B. Alur WhatsApp Tamu (Evolution API Webhook)
1. Tamu mengirim chat WhatsApp beruntun (misal 3 chat dalam 2 detik).
2. Evolution API menembakkan webhook `POST /webhook` ke `WebhookHandler`.
3. `WebhookHandler` mengekstrak pesan melalui `whatsapp.ExtractIncoming` dan memasukkannya ke `DebounceBuffer`.
4. `DebounceBuffer` menampung pesan per nomor HP dan me-reset timer debounce (4 detik).
5. Setelah timer berakhir tanpa ada pesan baru, seluruh pesan digabungkan menjadi satu pesan terpadu.
6. Buffer menjalankan pengecekan takeover dan memanggil Gemini Agent.
7. Respons akhir dipecah menjadi beberapa *bubble chat* alami yang dikirim satu per satu ke nomor tamu dengan delay simulasi mengetik manusia (*typing indicator*).
