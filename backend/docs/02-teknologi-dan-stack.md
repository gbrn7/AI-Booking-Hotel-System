# Modul 02: Teknologi & Stack Dependensi Backend

Backend sistem **Grand Oasis Hotel** mengadopsi stack teknologi modern yang berfokus pada performa eksekusi, efisiensi konsumsi memori (*low memory footprint*), dan keandalan konkurensi.

---

## 1. Daftar Teknologi Utama

### 1.1 Bahasa Pemrograman Go (Golang 1.25+)
* **Keunggulan Utama**:
  - Dikompilasi menjadi *single static binary* tanpa ketergantungan runtime eksternal.
  - Model konkurensi berbasis **Goroutine** dan **Channel** yang sangat ringan (~2KB per goroutine dibanding ~1MB pada thread OS konvensional), memungkinkan backend melayani ribuan sesi chat dan streaming SSE simultan dengan penggunaan RAM di bawah 50MB.
  - Dukungan kuat untuk cross-compilation (`CGO_ENABLED=0 GOOS=linux`) memudahkan kontainerisasi minimalis.

### 1.2 Gin Web Framework (`github.com/gin-gonic/gin`)
* **Peran**: Engine HTTP router dan middleware pipeline.
* **Fitur yang Dimanfaatkan**:
  - Radix tree router berkecepatan tinggi dengan alokasi memori minimal.
  - Penanganan streaming HTTP / Server-Sent Events (SSE) via `c.Stream()` dan `http.Flusher`.
  - JSON binding dan validasi struct otomatis (`binding:"required"`).
  - Middleware CORS kustom untuk komunikasi lintas origin dengan Next.js frontend.

### 1.3 GORM ORM (`gorm.io/gorm`)
* **Peran**: Object-Relational Mapping dan abstraksi database.
* **Fitur yang Dimanfaatkan**:
  - **Auto-Migration**: Pembuatan dan penyelarasan skema database secara otomatis pada startup tanpa perlu file migrasi SQL manual terpisah.
  - **Relational Preloading**: Preload relasi foreign key (`Preload("Room")`, `Preload("Guest")`) untuk efisiensi kueri.
  - **Pessimistic & Optimistic Concurrency**: Dukungan transaksi `db.Transaction()` untuk atomic reservation locking.

### 1.4 Dual Database Engine (PostgreSQL & SQLite)
* **PostgreSQL 16 (`gorm.io/driver/postgres`)**:
  - Engine utama untuk lingkungan production dan container Docker/Podman.
  - Keandalan transaksi ACID penuh, connection pooling (`pgx/v5`), dan isolasi konkurensi tinggi.
* **Pure-Go SQLite (`github.com/glebarez/sqlite`)**:
  - Driver SQLite 100% pure Go tanpa memerlukan CGO/GCC (`modernc.org/sqlite` backend).
  - Digunakan untuk unit testing otomatis berbasis in-memory (`:memory:`) yang sangat cepat dan simulator CLI offline.

### 1.5 Google GenAI SDK (`google.golang.org/genai`)
* **Peran**: SDK resmi Google untuk integrasi model bahasa Gemini (seperti `gemini-2.5-flash` dan `gemini-1.5-flash`).
* **Fitur yang Dimanfaatkan**:
  - Deklarasi skema tool menggunakan `FunctionDeclaration` berstandar JSON Schema.
  - Pemanggilan fungsi multi-langkah (*multi-turn tool calling loop*).
  - Konfigurasi parameter suhu (*temperature*), instruksi sistem (*system instructions*), dan handling token output.

### 1.6 Keamanan & Token (`golang.org/x/crypto/bcrypt` & `github.com/golang-jwt/jwt/v5`)
* **Bcrypt**: Pengacakan password staff dan admin hotel dengan cost factor 10, mencegah serangan rainbow table.
* **JWT (JSON Web Token)**: Penerbitan token HMAC-SHA256 untuk sesi login staf pada dashboard resepsionis, yang terverifikasi secara stateless di backend.

---

## 2. Tabel Ringkasan Dependensi (`go.mod`)

| Modul | Versi | Tujuan |
|---|---|---|
| `google.golang.org/genai` | `v1.73.0` | Client resmi Google Gemini AI & function calling |
| `github.com/gin-gonic/gin` | `v1.12.0` | HTTP Web Framework berkecepatan tinggi |
| `gorm.io/gorm` | `v1.31.2` | ORM untuk abstraksi database PostgreSQL & SQLite |
| `gorm.io/driver/postgres` | `v1.6.3` | Driver PostgreSQL GORM berkinerja tinggi |
| `github.com/glebarez/sqlite` | `v1.11.0` | Driver SQLite pure-Go tanpa dependensi CGO |
| `github.com/golang-jwt/jwt/v5`| `v5.3.1` | Pembuatan dan verifikasi token autentikasi JWT |
| `golang.org/x/crypto` | `v0.58.0` | Algoritma hashing password bcrypt |
| `github.com/joho/godotenv` | `v1.5.1` | Pemuatan file konfigurasi `.env` otomatis |
