# Modul 04: Live Web Chat Widget & SSE Streaming

## 1. Arsitektur Widget Chat Tamu

Komponen `LiveChatWidget.tsx` adalah jembatan interaksi langsung antara tamu website dan agen AI Rama. Widget ini dirancang sebagai komponen melayang (*floating modal*) di pojok kanan bawah yang tidak mengganggu penjelajahan halaman.

```
       Browser Tamu (Next.js)                            Go Backend Engine
                 │                                               │
          [Input Pesan]                                          │
                 │                                               │
                 ├────────── POST /api/chat/stream ─────────────►│
                 │           (session_id, message)               │
                 │                                               │
          [Mulai Stream]                                         │
                 │◄───────── HTTP 200 chunked SSE ───────────────┤
                 │                                               │
                 │◄── data: {"type": "bubble", "content": "..."} ┤
                 │                                               │
          (Render Bubble 1)                                      │
                 │                                               │
                 │◄── data: {"type": "bubble", "content": "..."} ┤
                 │                                               │
          (Render Bubble 2)                                      │
                 │                                               │
                 │◄── data: {"type": "done"} ────────────────────┤
                 │                                               │
          (Stream Selesai)                                       │
```

---

## 2. Cara Kerja Streaming SSE di Browser

Berbeda dengan WebSocket yang bersifat dua arah (*bidirectional*), **Server-Sent Events (SSE)** menggunakan protokol HTTP standar satu arah dari server ke klien yang lebih ringan, tahan firewall, dan tidak memerlukan server koneksi terpisah.

### Implementasi Pembaca Stream (`LiveChatWidget.tsx`):
```typescript
const response = await fetch(`${BACKEND_URL}/api/chat/stream`, {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({ session_id: sessionId, message: text }),
});

const reader = response.body.getReader();
const decoder = new TextDecoder("utf-8");
let buffer = "";

while (true) {
  const { done, value } = await reader.read();
  if (done) break;

  buffer += decoder.decode(value, { stream: true });
  const lines = buffer.split("\n\n");
  buffer = lines.pop() || "";

  for (const line of lines) {
    if (line.startsWith("data: ")) {
      const jsonStr = line.replace("data: ", "").trim();
      const parsed = JSON.parse(jsonStr);

      if (parsed.type === "bubble" || parsed.type === "takeover") {
        setMessages((prev) => [
          ...prev,
          {
            id: "asst-" + Date.now(),
            sender: parsed.type === "takeover" ? "human_agent" : "assistant",
            content: parsed.content,
            timestamp: new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }),
          },
        ]);
      }
    }
  }
}
```

---

## 3. Fitur Kenyamanan Pengguna (UX Enhancements)

### 3.1 Quick Chips (`QuickChips.tsx`)
Tamu dapat langsung mengeklik saran pertanyaan populer tanpa perlu mengetik manual:
* *"Cek kamar kosong akhir pekan ini"*
* *"Rekomendasi kamar untuk honeymoon"*
* *"Tarif Deluxe Ocean View berapa?"*
* *"Cara pesan dan batas hold booking"*

### 3.2 Persistensi Sesi (`localStorage`)
* Sesi percakapan tamu disimpan pada browser: `localStorage.getItem("grand_oasis_chat_sid")`.
* Format session ID: `web:<random_uuid>`.
* Tamu yang me-refresh browser atau berpindah halaman tetap dapat melihat riwayat obrolan sebelumnya via endpoint `GET /api/chat/history?session_id=...`.
* Tersedia tombol reset sesi (*New Chat*) untuk membersihkan riwayat obrolan jika tamu ingin memulai pertanyaan baru.

### 3.3 Integrasi Lintas Komponen (`open-ai-chat` Event)
Pada katalog kartu kamar (`RoomCard.tsx`), tombol **"Tanya Rama"** langsung memicu modal chat terbuka dan mengirim pertanyaan spesifik:
```typescript
window.dispatchEvent(
  new CustomEvent("open-ai-chat", {
    detail: { query: `Halo Rama, bisakah jelaskan detail dan promo untuk kamar ${room.name}?` },
  })
);
```
