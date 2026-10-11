# 📚 Modul Dokumentasi Teknis Frontend (Next.js 15 Fullstack)

Selamat datang di dokumentasi teknis mendalam untuk **Next.js 15 Frontend Portal & Admin Dashboard** pada sistem **Grand Oasis Hotel Booking Platform**.

Frontend ini dibangun menggunakan teknologi mutakhir: **Next.js 15 (App Router)**, **React 19**, **Tailwind CSS**, **Lucide Icons**, dan **NextAuth.js**, mencakup pengalaman publik tamu hotel serta pusat komando staf resepsionis (*Command Center*).

---

## 🗺️ Daftar Modul Pembelajaran & Referensi Teknis

| Modul | Topik | Deskripsi Utama |
|---|---|---|
| [**Modul 01**](./01-arsitektur-nextjs.md) | **Arsitektur Next.js 15 App Router** | Pemisahan Server Components & Client Components, Route Grouping, Layout Hierarchy, dan Rendering Strategy. |
| [**Modul 02**](./02-teknologi-dan-styling.md) | **Teknologi, Styling, & Desain Resor** | Tailwind CSS theme tokens (Slate & Gold Luxury Hospitality), tipografi serif, icon set Lucide React, dan standalone container output. |
| [**Modul 03**](./03-flow-data-dan-api-client.md) | **Flow Data & API Client Integration** | Library klien `lib/api.ts`, komunikasi dengan Go backend, pemetaan tipe data TypeScript, dan fallback state. |
| [**Modul 04**](./04-live-chat-widget-dan-sse.md) | **Live Web Chat Widget & SSE Streaming** | Arsitektur widget chat melayang (*floating widget*), parsing stream Server-Sent Events (SSE), Quick Chips, dan sinkronisasi sesi. |
| [**Modul 05**](./05-autentikasi-dan-middleware.md) | **Autentikasi NextAuth.js & Middleware Security** | Credentials Provider JWT, sesi otorisasi role `admin` & `receptionist`, serta proteksi route `/admin/*` via Next.js Middleware. |
| [**Modul 06**](./06-admin-command-center.md) | **Admin Command Center & Takeover Center** | Tampilan inbox dua kolom, saklar AI Takeover (ON/OFF), verifikasi bukti transfer booking HOLD, dan CMS inventori kamar. |

---

## 📂 Struktur Direktori Frontend

```
frontend/
├── src/
│   ├── app/
│   │   ├── (auth)/login/     # Halaman Login Staf & Kredensial
│   │   ├── admin/            # Area Terproteksi Admin & Resepsionis
│   │   │   ├── layout.tsx    # Shell Dashboard (Sidebar + Header)
│   │   │   ├── page.tsx      # Overview KPI Metrik & Shortcut
│   │   │   ├── takeover/     # Live Chat Takeover Center
│   │   │   ├── bookings/     # Tabel Reservasi & Verifikasi Pembayaran
│   │   │   └── rooms/        # Manajemen Kapasitas & Tarif Kamar
│   │   ├── api/auth/         # Route Handler NextAuth.js API
│   │   ├── layout.tsx        # Shell Utama Aplikasi & Floating Chat Widget
│   │   ├── page.tsx          # Landing Page Publik Hotel
│   │   └── globals.css       # Definisi Tailwind CSS & CSS Variables
│   ├── components/
│   │   ├── landing/          # Navbar, HeroSection, QuickSearch, RoomGrid, Amenities, Footer
│   │   ├── chat/             # LiveChatWidget, ChatMessage, QuickChips
│   │   └── admin/            # AdminSidebar, TakeoverChat, BookingVerificationModal
│   ├── lib/
│   │   ├── api.ts            # Typed HTTP API Client ke Go Backend
│   │   └── auth.ts           # Konfigurasi NextAuth.js Options & Callbacks
│   └── types/
│       └── index.ts          # Kontrak Interface Data TypeScript
├── public/                   # Asset Statis (.gitkeep, favicon, images)
├── docs/                     # Dokumentasi Teknis & Modul Pembelajaran
├── Dockerfile                # Multi-stage Dockerfile Standalone Runner
├── next.config.mjs           # Konfigurasi Next.js (Standalone, Image domain)
├── package.json              # Dependensi npm & skrip build
├── tailwind.config.ts        # Kustomisasi warna tema & breakpoint
└── tsconfig.json             # Konfigurasi TypeScript Compiler
```
