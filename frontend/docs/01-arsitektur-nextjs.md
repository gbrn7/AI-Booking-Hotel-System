# Modul 01: Arsitektur Next.js 15 App Router

## 1. Paradigma Next.js 15 App Router

Frontend **Grand Oasis Hotel** mengadopsi model **Next.js 15 App Router** (`src/app/`) yang memadukan keunggulan performa Server-Side Rendering (SSR) untuk SEO publik dengan interaktivitas instan Client-Side Rendering (CSR) untuk live chat dan dashboard operasional.

### Diagram Hirarki Layout & Halaman
```
src/app/layout.tsx (Root Shell)
├── globals.css & Fonts
├── <LiveChatWidget /> (Floating Widget Global di Semua Halaman)
│
├── (Public Route)
│   └── page.tsx (Landing Page Hotel)
│       ├── <Navbar />
│       ├── <HeroSection /> (<QuickSearch />)
│       ├── <RoomGrid /> (<RoomCard />)
│       ├── <Amenities />
│       └── <Footer />
│
├── (Auth Route)
│   └── login/page.tsx (Form Login Kredensial Staf)
│
└── admin/layout.tsx (Dashboard Shell - Terproteksi Middleware)
    ├── <AdminSidebar /> (Navigasi Persisten)
    ├── Header Status Staf
    └── Page Content:
        ├── page.tsx (KPI Overview Metrik)
        ├── takeover/page.tsx (<TakeoverChat />)
        ├── bookings/page.tsx (<BookingVerificationModal />)
        └── rooms/page.tsx (CMS Inventori Kamar)
```

---

## 2. Server Components vs Client Components

Next.js 15 secara default memperlakukan seluruh komponen sebagai **React Server Components (RSC)** kecuali jika dideklarasikan dengan direktif `"use client"`.

### Strategi Pembagian Komponen:
1. **Server Components**:
   - `src/app/layout.tsx`, `src/app/page.tsx`, `src/app/admin/layout.tsx`.
   - Menghasilkan markup HTML statis yang sangat cepat di-load pertama kali oleh browser tanpa membebani bundle JavaScript (*Zero Bundle Size* untuk komponen statis).
   - Mengoptimalkan Search Engine Optimization (SEO) landing page hotel.
2. **Client Components (`"use client"`)**:
   - `QuickSearch.tsx`: Memerlukan state tanggal dan event handler pencarian.
   - `LiveChatWidget.tsx` & `ChatMessage.tsx`: Memerlukan koneksi streaming SSE, `localStorage` sesi, dan manipulasi DOM scroll otomatis.
   - `TakeoverChat.tsx`: Memerlukan interval polling real-time, state aktif takeover, dan input balasan manual staf.
   - `BookingVerificationModal.tsx`: Memerlukan state modal popup dan konfirmasi mutasi API.

---

## 3. Komposisi Floating Live Chat Global

Salah satu keunggulan arsitektur ini adalah penempatan `<LiveChatWidget />` langsung di dalam `src/app/layout.tsx`:
```tsx
export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="id">
      <body className="min-h-screen bg-slate-50 text-slate-900 antialiased relative">
        {children}
        <LiveChatWidget />
      </body>
    </html>
  );
}
```
* Dengan pola ini, widget chat asisten Rama tetap persisten dan terbuka ketika tamu menelusuri berbagai informasi di landing page tanpa me-reset riwayat percakapan.
* Komponen publik seperti tombol *"Tanya Rama"* pada `RoomCard.tsx` dapat memicu widget chat secara instan melalui Custom Event `window.dispatchEvent(new CustomEvent("open-ai-chat", ...))`.
