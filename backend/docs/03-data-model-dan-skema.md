# Modul 03: Data Model & Skema Database

Sistem **Grand Oasis Hotel** mengelola 6 entitas utama di dalam database untuk melacak kamar, data tamu, reservasi, histori percakapan multi-channel, dan otorisasi pengguna staf.

---

## 1. Entity-Relationship Diagram (ERD)

```
  ┌────────────────────────┐                   ┌────────────────────────┐
  │         rooms          │                   │         guests         │
  ├────────────────────────┤                   ├────────────────────────┤
  │ PK id                  │                   │ PK id                  │
  │    name                │                   │    name                │
  │    type                │                   │    phone_number (UQ)   │
  │    base_price_per_night│                   │    email               │
  │    capacity_guests     │                   └───────────┬────────────┘
  │    total_units         │                               │
  │    amenities (JSON)    │                               │ 1
  │    image_urls (JSON)   │                               │
  │    is_active           │                               │
  └───────────┬────────────┘                               │
              │ 1                                          │
              │                                            │
              │ N                                        N │
        ┌─────┴────────────────────────────────────────────┴─────┐
        │                        bookings                        │
        ├────────────────────────────────────────────────────────┤
        │ PK id                                                  │
        │    booking_code (UQ)                                   │
        │ FK guest_id ────────────► guests.id                    │
        │ FK room_id  ────────────► rooms.id                     │
        │    check_in_date                                       │
        │    check_out_date                                      │
        │    total_nights                                        │
        │    total_rooms                                         │
        │    total_price                                         │
        │    status (HOLD, CONFIRMED, CANCELLED, ...)            │
        │    hold_expires_at                                     │
        │    payment_proof_url                                   │
        │    special_requests                                    │
        │ FK confirmed_by_user_id ─► users.id (nullable)         │
        └────────────────────────────────────────────────────────┘

  ┌────────────────────────┐                   ┌────────────────────────┐
  │     conversations      │ 1               N │        messages        │
  ├────────────────────────┼───────────────────┼────────────────────────┤
  │ PK id (wa:../web:..)   │                   │ PK id                  │
  │ FK guest_id (nullable) │                   │ FK conversation_id     │
  │    channel             │                   │    sender              │
  │    is_human_takeover   │                   │    content             │
  │    last_message_at     │                   │    thought_process     │
  └────────────────────────┘                   └────────────────────────┘

  ┌────────────────────────┐
  │         users          │
  ├────────────────────────┤
  │ PK id                  │
  │    username (UQ)       │
  │    email (UQ)          │
  │    password_hash       │
  │    role (admin/staff)  │
  └────────────────────────┘
```

---

## 2. Rincian Entitas Data (`internal/model/models.go`)

### 2.1 Entitas `Room` (Kamar Hotel)
Mewakili tipe kamar fisik hotel dan jumlah kuota total yang dimiliki properti.
* `id` (`uint`, Primary Key)
* `name` (`string`, max 100): Nama pemasaran kamar (contoh: *"Deluxe Ocean View"*).
* `type` (`string`, index): Kategori teknis (`standard`, `deluxe`, `suite`, `family_suite`).
* `base_price_per_night` (`float64`): Tarif dasar kamar per malam dalam Rupiah (IDR).
* `capacity_guests` (`int`): Batas jumlah tamu standar dalam satu kamar.
* `total_units` (`int`): Total inventori unit fisik kamar tipe ini yang dimiliki hotel (misal: 8 unit).
* `description` (`string`, text): Penjelasan daya tarik, pemandangan, dan spesifikasi kamar.
* `amenities` (`string`, JSON text): Daftar fasilitas dalam bentuk array JSON string.
* `image_urls` (`string`, JSON text): URL foto-foto resolusi tinggi dari kamar.
* `is_active` (`bool`): Menandai apakah kamar dapat disewa atau sedang renovasi.

### 2.2 Entitas `Guest` (Tamu Hotel)
Menyimpan identitas tamu hotel yang diidentifikasi secara unik melalui nomor kontak/WhatsApp.
* `id` (`uint`, Primary Key)
* `name` (`string`): Nama lengkap tamu.
* `phone_number` (`string`, Unique Index): Kunci identitas utama lintas channel WhatsApp dan Web.
* `email` (`string`, opsional): Alamat surel tamu.

### 2.3 Entitas `Booking` (Reservasi Kamar)
Jantung operasional transaksi yang mencatat pemesanan kamar beserta rentang tanggal dan status pembayaran.
* `booking_code` (`string`, Unique): Format `HTL-YYMMDD-XXXXX` (contoh: `HTL-261011-8X9Y`).
* `check_in_date` & `check_out_date` (`string`, `YYYY-MM-DD`): Tanggal mulai dan selesai menginap.
* `total_nights` (`int`): Jumlah malam yang dihitung dari selisih kedua tanggal.
* `total_rooms` (`int`): Jumlah unit kamar yang direservasi (default: 1).
* `total_price` (`float64`): `total_nights * base_price_per_night * total_rooms`.
* `status` (`string`): Status siklus reservasi.
* `hold_expires_at` (`*time.Time`, index): Batas waktu kedaluwarsa status `HOLD` (default: 2 jam).

### 2.4 Entitas `Conversation` & `Message`
* `Conversation`: Mengidentifikasi sesi chat dengan ID gabungan channel (contoh: `wa:62812345678` atau `web:e8a9f2...`). Menyimpan flag `is_human_takeover`.
* `Message`: Catatan setiap bubble percakapan dengan peran sender (`user`, `assistant`, `agent_tool`, `human_agent`) dan *thought process* internal AI.

### 2.5 Entitas `User` (Staf & Admin)
Menyimpan akun kredensial staf hotel dengan role `admin` (General Manager) atau `receptionist` (Front Desk Staff).

---

## 3. State Machine Siklus Reservasi (Booking Status Lifecycle)

```
              ┌─────────────────────────┐
              │   AI/Web Buat Booking   │
              └────────────┬────────────┘
                           │
                           ▼
                 ┌───────────────────┐
                 │    Status HOLD    │◄─── Unit kamar terkunci sementara (2 jam)
                 └──────┬─────┬──────┘
       Waktu Habis      │     │  Staf Resepsionis
    (Background Worker) │     │  Verifikasi Pembayaran
                        │     │
                        ▼     ▼
     ┌────────────────────┐ ┌────────────────────────┐
     │  Status CANCELLED  │ │    Status CONFIRMED    │◄── Reservasi resmi lunas
     │ (Unit dilepas lagi)│ └───────────┬────────────┘
     └────────────────────┘             │ Tamu Tiba di Hotel
                                        ▼
                            ┌────────────────────────┐
                            │   Status CHECKED_IN    │
                            └───────────┬────────────┘
                                        │ Tamu Selesai Menginap
                                        ▼
                            ┌────────────────────────┐
                            │   Status CHECKED_OUT   │
                            └────────────────────────┘
```
