# Modul 03: Flow Data & Integrasi API Client

## 1. Arsitektur Komunikasi Frontend & Backend

Frontend Next.js berkomunikasi dengan Go backend (`http://localhost:8080`) melalui pustaka klien tersentralisasi pada `src/lib/api.ts`. Pola ini memisahkan logika antarmuka komponen dari rincian HTTP request (URL endpoint, header, dan serialisasi JSON).

```
  ┌─────────────────────────────────────────────────────────────┐
  │                 React UI Components                         │
  │ (RoomGrid, QuickSearch, LiveChatWidget, TakeoverChat, etc.) │
  └──────────────────────────────┬──────────────────────────────┘
                                 │
                                 ▼
  ┌─────────────────────────────────────────────────────────────┐
  │                 Typed API Client (lib/api.ts)               │
  │  - fetchRooms()                                             │
  │  - checkRoomAvailability(checkIn, checkOut)                 │
  │  - fetchAdminStats(token)                                   │
  │  - fetchAdminConversations(token)                           │
  │  - toggleTakeover(id, isTakeover, token)                    │
  │  - confirmBookingPayment(id, token)                         │
  └──────────────────────────────┬──────────────────────────────┘
                                 │ HTTP REST / SSE Stream
                                 ▼
  ┌─────────────────────────────────────────────────────────────┐
  │                   Go Gin Backend Engine                     │
  │                    (Port 8080 / DB GORM)                    │
  └─────────────────────────────────────────────────────────────┘
```

---

## 2. Definisi Kontrak Tipe Data (`src/types/index.ts`)

Seluruh respons dari Go backend dijamin ketepatan tipenya menggunakan antarmuka TypeScript yang selaras dengan GORM struct backend:

```typescript
export interface Room {
  id: number;
  name: string;
  type: string;
  base_price_per_night: number;
  capacity_guests: number;
  total_units: number;
  description: string;
  amenities: string;
  image_urls: string;
  is_active: boolean;
}

export interface Booking {
  id: number;
  booking_code: string;
  guest_id: number;
  guest?: Guest;
  room_id: number;
  room?: Room;
  check_in_date: string;
  check_out_date: string;
  total_nights: number;
  total_rooms: number;
  total_price: number;
  status: "HOLD" | "CONFIRMED" | "CANCELLED" | "CHECKED_IN" | "CHECKED_OUT";
  hold_expires_at?: string;
}
```

---

## 3. Flow Data Utama

### 3.1 Pencarian Ketersediaan Cepat (*Quick Search*)
1. Tamu memilih tanggal check-in, check-out, dan jumlah tamu di `QuickSearch.tsx`.
2. Komponen memanggil `checkRoomAvailability(checkIn, checkOut)`:
   ```typescript
   export async function checkRoomAvailability(checkIn: string, checkOut: string, roomType?: string) {
     const params = new URLSearchParams({ check_in: checkIn, check_out: checkOut });
     const res = await fetch(`${BACKEND_URL}/api/availability?${params.toString()}`, { cache: "no-store" });
     const data = await res.json();
     return data.availability || [];
   }
   ```
3. Backend mengevaluasi inventori per malam dan mengembalikan sisa unit riil (`available_units`).
4. Hasil ditampilkan langsung ke tamu, dan halaman melakukan *smooth scroll* ke katalog kamar.

### 3.2 Polling Percakapan di Dashboard Resepsionis
1. Komponen `TakeoverChat.tsx` memuat daftar obrolan aktif melalui `fetchAdminConversations()`.
2. Menggunakan `setInterval` (interval 5 detik), komponen memperbarui status pesan masuk dan status takeover secara otomatis tanpa perlu me-refresh halaman web secara manual.
