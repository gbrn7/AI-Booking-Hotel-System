# Modul 06: Receptionist Command Center & Takeover UI

## 1. Konsep Receptionist Command Center

Dashboard Admin **Grand Oasis Hotel** (`/admin`) dirancang khusus sebagai **Pusat Komando Staf Resepsionis** (*Front Desk Command Center*). Tujuannya adalah memberikan transparansi penuh atas interaksi AI Sales Agent dengan tamu, memudahkan verifikasi mutasi bank untuk pembayaran kamar, dan menyediakan tombol kendali darurat intervensi manusia (*Human Takeover*).

---

## 2. Fitur-Fitur Utama Command Center

### 2.1 Live Takeover Center (`/admin/takeover`)
* **Tampilan Dua Kolom**:
  - **Kolom Kiri**: Daftar seluruh percakapan aktif dari saluran WhatsApp maupun Web Live Chat, dilengkapi badge saluran dan indikator waktu pesan terakhir.
  - **Kolom Kanan**: Transkrip percakapan real-time dengan pembedaan avatar jelas antara Tamu, Rama (AI), dan Staf Resepsionis.
* **Saklar Takeover (AI ON / OFF)**:
  - Tombol toggle untuk mengalihkan kendali obrolan ke manusia.
  - Ketika Takeover aktif, backend menonaktifkan balasan Gemini, dan staf front desk dapat mengetik balasan manual langsung dari browser yang otomatis terkirim ke WhatsApp / Web tamu.

### 2.2 Manajemen Reservasi & Verifikasi Pembayaran (`/admin/bookings`)
* **Filter Tab Status**:
  - `ALL`: Seluruh reservasi.
  - `HOLD`: Reservasi sementara yang menunggu transfer uang dan verifikasi staf (batas waktu 2 jam).
  - `CONFIRMED`: Reservasi yang telah lunas dan diverifikasi.
  - `CANCELLED`: Reservasi yang dibatalkan tamu atau kedaluwarsa oleh background worker.
* **Modal Verifikasi Pembayaran (`BookingVerificationModal.tsx`)**:
  - Menampilkan ringkasan kode booking, nama tamu, kamar, durasi, dan total nominal tagihan.
  - Menampilkan rekening resmi penerima: **BCA 8870-1234-5678-9000 a/n Grand Oasis Hotel PT**.
  - Tombol aksi **"Verifikasi & Set CONFIRMED"** yang memicu endpoint `POST /api/admin/bookings/:id/confirm` dan mengunci status kamar secara permanen.

### 2.3 CMS Inventori Kamar (`/admin/rooms`)
* Menampilkan kartu seluruh tipe kamar yang terdaftar di sistem.
* Menampilkan jumlah unit fisik yang dimiliki hotel (`total_units`), tarif dasar per malam, batas kapasitas tamu, dan status aktif.
* Angka total unit ini yang digunakan oleh Go inventory engine untuk kalkulasi matematis kapasitas per malam.

### 2.4 Ringkasan KPI Overview (`/admin`)
* Kartu statistik metrik utama hotel:
  1. **Total Tipe Kamar** aktif dalam katalog.
  2. **Antrean HOLD**: Jumlah reservasi yang mendesak untuk diverifikasi.
  3. **Total Reservasi CONFIRMED**: Jumlah pesanan sah yang telah lunas.
  4. **Percakapan Masuk**: Total sesi chat tamu aktif di sistem.
