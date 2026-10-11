# Modul 02: Teknologi, Styling, & Desain Resor Mewah

## 1. Stack Teknologi Frontend

| Teknologi | Versi | Peran Utama |
|---|---|---|
| **Next.js** | `v15.1+` | Framework React Fullstack dengan App Router dan optimasi bundle. |
| **React** | `v19.0+` | Pustaka antarmuka deklaratif dengan dukungan Concurrent Mode. |
| **Tailwind CSS** | `v3.4+` | Utility-first CSS framework untuk styling responsif berkecepatan tinggi. |
| **Lucide React** | `v0.468+` | Kumpulan ikon SVG konsisten untuk indikator perhotelan dan aksi UI. |
| **NextAuth.js** | `v4.24+` | Solusi autentikasi fleksibel untuk mengamankan area dashboard staf. |

---

## 2. Sistem Desain: *Luxury Beachfront Sanctuary*

Desain antarmuka **Grand Oasis Hotel** mengusung tema estetika resor bintang 5 modern: hangat, tenang, dan terstruktur rapi.

### 2.1 Palet Warna Tema (`tailwind.config.ts`)
* **Warm Amber & Gold (`colors.brand`)**:
  - `brand-500: #d99324`, `brand-400: #e3ad44`, `brand-50: #fdf8ee`.
  - Digunakan untuk aksen premium, bintang hotel, tombol CTA utama, badge concierge AI, dan sorotan harga kamar.
* **Deep Navy & Slate (`colors.slate` & `colors.navy`)**:
  - `slate-900: #0f172a`, `slate-800: #1e293b`.
  - Memberikan nuansa elegan dan keterbacaan tinggi pada teks judul, header, dan sidebar dashboard admin.
* **Emerald & Rose Semantic Colors**:
  - `emerald-600`: Melambangkan status `CONFIRMED`, indikator online concierge, dan tombol verifikasi pembayaran.
  - `amber-600`: Melambangkan status `HOLD` (menunggu konfirmasi).
  - `rose-600`: Melambangkan status pembatalan atau saklar mode `Takeover ON`.

### 2.2 Tipografi
* **Font Serif (`font-serif`)**: Digunakan untuk judul utama (headings), nama hotel (*Grand Oasis*), dan nominal harga kamar untuk menonjolkan karakter klasik mewah.
* **Font Sans-Serif (`font-sans`)**: Digunakan untuk paragraf deskripsi, form input, dan teks chat agar mudah dibaca di layar ponsel maupun desktop.

---

## 3. Optimasi Container Standalone (`output: "standalone"`)

Pada file `next.config.mjs`:
```javascript
const nextConfig = {
  output: "standalone",
  images: {
    remotePatterns: [{ protocol: "https", hostname: "images.unsplash.com" }],
  },
};
export default nextConfig;
```

### Mengapa Fitur Standalone Sangat Penting?
1. Saat `npm run build` dijalankan, Next.js menganalisis pohon ketergantungan dan mengekstrak **hanya file-file yang benar-benar dibutuhkan** ke folder `.next/standalone`.
2. Image Docker tidak perlu menyertakan seluruh folder `node_modules` pengembangan (~500MB).
3. Image runtime `frontend/Dockerfile` hanya berukuran **sekitar ~100MB** dan dijalankan cukup dengan binary Node.js bawaan:
   ```dockerfile
   CMD ["node", "server.js"]
   ```
