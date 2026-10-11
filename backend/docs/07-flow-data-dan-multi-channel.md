# Modul 07: Flow Data & Multi-Channel Gateway (WhatsApp & Web SSE)

## 1. Arsitektur Multi-Channel

Backend **Grand Oasis Hotel** dirancang untuk melayani dua saluran komunikasi tamu secara simultan menggunakan logika bisnis yang sama:
1. **WhatsApp Channel**: Melalui gateway **Evolution API** dengan integrasi webhook.
2. **Web Live Chat Channel**: Melalui **Server-Sent Events (SSE)** yang tertanam pada website Next.js.

```
       WhatsApp Tamu                              Browser Tamu (Web)
             │                                             │
      (Chat Beruntun)                               (Input Pesan Web)
             │                                             │
             ▼                                             ▼
    Evolution API Gateway                         POST /api/chat/stream
             │                                             │
             ▼                                             │
     POST /webhook                                         │
             │                                             │
             ▼                                             │
     Debounce Buffer (4s)                                  │
             │                                             │
             └──────────────────────┬──────────────────────┘
                                    │
                                    ▼
                         Cek Status Takeover?
                         ┌──────────┴──────────┐
                  (Aktif)│                     │(Nonaktif)
                         ▼                     ▼
             Simpan ke Inbox Staf     Panggil Gemini AI
             (Tanpa Respons AI)                │
                                               ▼
                                      Respon Bubble [NEXT]
                                               │
                                      ┌────────┴────────┐
                                      ▼                 ▼
                              Kirim WhatsApp      Stream via SSE
                              (Evolution API)     (Chunked Browser)
```

---

## 2. Mengapa Perlu Debounce Buffer di WhatsApp?

### Masalah Pengguna WhatsApp:
Pengguna WhatsApp hampir tidak pernah mengetik pesan panjang dalam satu kali kirim. Mereka biasanya mengirim pesan berturut-turut:
* *10:00:01* — "Halo"
* *10:00:02* — "Mau tanya kamar untuk weekend ini"
* *10:00:04* — "Ada yang deluxe ga ya?"

Jika setiap webhook langsung dikirim ke AI:
1. AI akan merespons 3 kali secara terpisah.
2. Respons pertama hanya akan menjawab "Halo", respons kedua bingung konteks tanggalnya, respons ketiga baru memproses kamar.
3. Terjadi pemborosan token API dan pengalaman pengguna yang buruk.

### Solusi Debounce Buffer (`internal/buffer/buffer.go`):
1. Setiap pesan masuk disimpan ke antrean memori map per nomor telepon: `buffers[phone]`.
2. Jika sudah ada timer aktif untuk nomor tersebut, timer dibatalkan (`t.Stop()`).
3. Timer baru berdurasi 4 detik dibuat (`time.AfterFunc`).
4. Jika tidak ada pesan baru selama 4 detik, seluruh pesan digabung dengan pemisah baris baru (`\n`) dan diproses satu kali saja oleh AI.

---

## 3. Web Chat SSE Streaming (`internal/handler/chat.go`)

Berbeda dengan WhatsApp yang bersifat asinkron melalui webhook, tamu website menginginkan pengalaman respons cepat (*low latency*).

Oleh karena itu, endpoint `/api/chat/stream` menggunakan **Server-Sent Events (SSE)**:
* Header respons:
  ```http
  Content-Type: text/event-stream
  Cache-Control: no-cache
  Connection: keep-alive
  Transfer-Encoding: chunked
  ```
* Format event streaming:
  ```
  data: {"type": "bubble", "index": 0, "content": "Halo Kak! Selamat datang di Grand Oasis Hotel 🌊"}

  data: {"type": "bubble", "index": 1, "content": "Untuk tanggal tersebut, kamar Deluxe Ocean View kami masih tersedia."}

  data: {"type": "done"}
  ```
* Frontend Next.js membaca stream menggunakan `response.body.getReader()` dan merender bubble chat secara dinamis dengan animasi halus.

---

## 4. Mekanisme Human Takeover (Intervensi Resepsionis)

Ketika tamu mengajukan komplain, meminta diskon khusus, atau tool `transfer_to_receptionist` dipanggil:
1. Kolom `is_human_takeover` di tabel `conversations` diubah menjadi `true`.
2. Seluruh pesan baru dari tamu tersebut tetap dicatat di database, namun **pemanggilan AI dibatalkan**.
3. Notifikasi badge merah "Takeover" menyala di dashboard resepsionis Next.js.
4. Resepsionis mengetik balasan manual melalui endpoint `POST /api/admin/conversations/:id/reply`.
5. Jika channel tamu adalah WhatsApp, balasan staf diteruskan ke Evolution API. Jika Web Chat, balasan tersimpan dan dapat dibaca tamu.
6. Resepsionis dapat menekan tombol **"Kembalikan ke AI"** kapan saja untuk mengaktifkan kembali respons otomatis Rama Concierge.
