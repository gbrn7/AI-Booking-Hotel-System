# Modul 06: Gemini Agentic Loop & 6 Tools Bisnis Hotel

## 1. Konsep Agentic Loop vs Prompt Statis

Pendekatan kecerdasan buatan konvensional seringkali hanya mengirimkan satu prompt teks ke LLM dan berharap LLM memberikan jawaban tanpa akses data riil. Hal ini menyebabkan **halusinasi harga dan ketersediaan kamar**.

Pada sistem **Grand Oasis Hotel**, kita mengimplementasikan **Autonomous Agentic Loop** menggunakan **Google GenAI SDK**:
* Model Gemini bertindak sebagai **otak pengambil keputusan**.
* Model memiliki akses ke **6 Tools Bisnis Hotel** yang terhubung langsung ke database PostgreSQL/SQLite.
* Model dapat memanggil tool secara otonom berkali-kali (*multi-turn function calling*) sebelum merangkum jawaban akhir kepada tamu.

---

## 2. Arsitektur Loop Eksekusi (`internal/agent/agent.go`)

```
   Pesan Tamu + Konteks Waktu & No HP + Riwayat Percakapan
                         │
                         ▼
             ┌───────────────────────┐
             │ a.client.Models.      │◄─────────────────────────────┐
             │ GenerateContent()     │                              │
             └───────────┬───────────┘                              │
                         │                                          │
                         ▼                                          │
                 Ada Function Call?                                 │
                 ┌───────┴───────┐                                  │
                 ▼               ▼                                  │
               (TIDAK)          (YA)                                │
                 │               │                                  │
                 │               ▼                                  │
                 │    Eksekusi Tools di Go:                         │
                 │    - a.tools.ExecuteTool(name, args)             │
                 │               │                                  │
                 │               ▼                                  │
                 │    Kirim FunctionResponse balik ke percakapan ───┘
                 │    (Role: "user", Parts: FunctionResponse)
                 │    (Maksimal 5 iterasi)
                 │
                 ▼
    Sintesis Balasan Teks Terakhir
                 │
                 ▼
    Pemisahan Bubble Chat via [NEXT]
```

---

## 3. Rincian 6 Tools Bisnis Hotel (`internal/tools/tools.go`)

| No | Nama Tool | Parameter | Fungsi Bisnis |
|---|---|---|---|
| 1 | `search_rooms` | `query` (kata kunci), `room_type`, `max_price` | Mencari tipe kamar berdasarkan kriteria tamu (fasilitas, balkon, bathtub, pemandangan, budget). Wajib dipanggil sebelum menyebut harga kamar. |
| 2 | `check_room_availability` | `check_in_date` (*YYYY-MM-DD*), `check_out_date` (*YYYY-MM-DD*), `room_type` | Menghitung sisa unit kamar kosong pada rentang tanggal menginap menggunakan engine kapasitas atomik. |
| 3 | `hold_booking` | `guest_name`, `guest_phone`, `room_id`, `check_in_date`, `check_out_date`, `total_rooms`, `special_requests` | Mengunci unit kamar ke status `HOLD` (2 jam), menghasilkan kode booking unik (`HTL-XXXXX`), dan menghitung total tagihan serta info rekening resmi. |
| 4 | `get_guest_bookings` | `phone`, `booking_code` | Melacak status pesanan aktif atau histori booking milik tamu. |
| 5 | `cancel_booking` | `booking_code`, `phone`, `reason` | Membatalkan reservasi yang masih berstatus `HOLD` dan melepaskan unit kembali ke inventori. |
| 6 | `transfer_to_receptionist`| `phone`, `reason` | Mengaktifkan saklar `is_human_takeover = true` pada sesi obrolan saat ada komplain atau negosiasi khusus, menghentikan respons otomatis AI. |

---

## 4. Persona Rama & Formatting Bubble Split (`[NEXT]`)

Di dalam template instruksi sistem (`backend/prompts/system.txt`):
* **Hospitality Tone**: Menggunakan bahasa Indonesia yang elegan, sopan, dan hangat layaknya concierge hotel bintang 5.
* **Bubble Splitter (`[NEXT]`)**: Tamu WhatsApp dan Web tidak menyukai membaca dinding teks panjang. Agen dilatih memisahkan responsnya menjadi beberapa bubble alami dengan tag `[NEXT]`:
  ```
  Halo Kak! Selamat datang di Grand Oasis Hotel 🌊[NEXT]
  Untuk tanggal 15-17 Oktober, kamar Deluxe Ocean View kami masih tersedia 2 unit dengan tarif Rp 1.250.000/malam.[NEXT]
  Apakah Kakak ingin Rama bantu buatkan hold reservasi sekarang? 😊
  ```
* Fungsi `whatsapp.SplitBubbles(reply)` di Go secara cerdas memotong teks ini menjadi array string untuk dikirimkan secara bertahap dengan animasi mengetik.
