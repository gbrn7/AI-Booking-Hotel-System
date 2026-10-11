# Modul 08: Kontrak Admin REST API & CLI Simulator

## 1. Spesifikasi Kontrak REST API

Semua endpoint backend disajikan pada port `8080` (dapat diubah via variabel `PORT`).

### 1.1 Endpoint Publik (Tamu & Web Portal)

#### `GET /api/rooms`
* **Deskripsi**: Mengambil daftar seluruh kamar hotel yang aktif beserta harga dan fasilitas.
* **Response `200 OK`**:
  ```json
  {
    "rooms": [
      {
        "id": 1,
        "name": "Superior King Room",
        "type": "standard",
        "base_price_per_night": 750000,
        "capacity_guests": 2,
        "total_units": 10,
        "description": "Kamar nyaman dengan pemandangan kota...",
        "amenities": "[\"King Bed\", \"Free WiFi\"]",
        "is_active": true
      }
    ]
  }
  ```

#### `GET /api/availability`
* **Query Params**: `check_in` (*YYYY-MM-DD*), `check_out` (*YYYY-MM-DD*), `room_type` (*opsional*).
* **Deskripsi**: Mengecek sisa unit kamar yang kosong pada rentang tanggal tersebut.
* **Response `200 OK`**:
  ```json
  {
    "availability": [
      {
        "room_id": 2,
        "name": "Deluxe Ocean View",
        "available_units": 6,
        "total_nights": 2,
        "estimated_total_price": 2500000
      }
    ]
  }
  ```

#### `POST /api/chat/stream`
* **Headers**: `Content-Type: application/json`
* **Request Body**:
  ```json
  {
    "session_id": "web:guest_123",
    "message": "Apakah kamar deluxe masih ada untuk tanggal 20-22 Oktober?",
    "phone": "08123456789"
  }
  ```
* **Response**: Stream Server-Sent Events (`text/event-stream`).

---

### 1.2 Endpoint Admin & Resepsionis (`/api/admin/*`)

#### `POST /api/admin/login`
* **Request Body**:
  ```json
  {
    "username": "admin",
    "password": "admin123"
  }
  ```
* **Response `200 OK`**: Mengembalikan token JWT dan data profil pengguna.

#### `GET /api/admin/stats`
* **Deskripsi**: Ringkasan metrik KPI untuk dashboard (total kamar, pending hold, confirmed, chat aktif).

#### `GET /api/admin/conversations`
* **Deskripsi**: Mengambil daftar seluruh sesi obrolan tamu dari channel WhatsApp dan Web.

#### `POST /api/admin/conversations/:id/takeover`
* **Request Body**: `{"is_takeover": true}` atau `{"is_takeover": false}`.
* **Deskripsi**: Menyalakan atau mematikan saklar intervensi manual staf hotel.

#### `POST /api/admin/conversations/:id/reply`
* **Request Body**: `{"message": "Halo Kak, saya staf front desk. Ada yang bisa saya bantu?"}`
* **Deskripsi**: Mengirimkan pesan manual dari resepsionis langsung ke tamu.

#### `GET /api/admin/bookings`
* **Query Params**: `status` (*HOLD*, *CONFIRMED*, *CANCELLED*).
* **Deskripsi**: Daftar reservasi kamar dengan rincian tamu, tanggal, total harga, dan status.

#### `POST /api/admin/bookings/:id/confirm`
* **Deskripsi**: Memverifikasi pembayaran dan mengubah status booking dari `HOLD` menjadi `CONFIRMED`.

---

## 2. Interactive CLI Simulator (`cmd/cli/main.go`)

Simulator terminal interaktif memungkinkan pengembang menguji alur percakapan AI Sales Agent, logika pemanggilan tools, dan pembuatan booking tanpa memerlukan WhatsApp atau browser.

### Cara Menjalankan
```bash
cd backend
go run ./cmd/cli
```

### Fitur Simulator
1. **Mode Offline/Mock**: Jika `GEMINI_API_KEY` belum diset, simulator otomatis berjalan dalam mode mock tanpa error crash.
2. **Eksekusi Tools Transparan**: Menampilkan log `🔍 [Internal Tool Execution]` saat AI memanggil function call di latar belakang.
3. **Format Bubble Realistis**: Balasan ditampilkan per bubble chat persis seperti tampilan di WhatsApp dan widget web.
4. **Perintah Keluar**: Cukup ketik `exit` atau `quit` untuk mengakhiri simulasi.
