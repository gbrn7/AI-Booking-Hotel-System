# Modul 05: Background Worker Pembersih Hold Expiry

## 1. Urgensi Bisnis Hold Expiration

Dalam bisnis reservasi hotel, tamu yang membuat pemesanan sementara (*HOLD Booking*) tetapi tidak kunjung melakukan transfer pembayaran akan menyebabkan **Inventory Starvation (Ghost Bookings)**. Kamar tampak penuh di sistem, tamu potensial lain tidak bisa memesan, namun hotel tidak mendapatkan pemasukan.

Oleh karena itu, sistem menerapkan kebijakan:
* Setiap reservasi baru otomatis diberi status `HOLD` dengan batas waktu `hold_expires_at = NOW() + HOLD_DURATION_MINUTES` (default: 2 jam atau 120 menit).
* Jika dalam 2 jam staf resepsionis belum memverifikasi pembayaran (status belum berubah menjadi `CONFIRMED`), reservasi harus dibatalkan secara otomatis agar kuota kamar langsung kembali ke inventori.

---

## 2. Arsitektur Background Worker Go (`internal/service/expiry.go`)

### Diagram Siklus Worker
```
       Server Startup
             │
             ▼
┌─────────────────────────┐
│ StartHoldExpiryWorker() │
└────────────┬────────────┘
             │ (spawn goroutine)
             ▼
      ┌─────────────┐
      │ time.Ticker │ ◄─── Berjalan setiap EXPIRY_CHECK_INTERVAL_SECONDS (default 60s)
      └──────┬──────┘
             │ (setiap detak ticker)
             ▼
┌──────────────────────────────┐
│  ExpireStaleHoldBookings()   │
│  UPDATE bookings             │
│  SET status = 'CANCELLED'    │
│  WHERE status = 'HOLD'       │
│    AND hold_expires_at < NOW()│
└────────────┬─────────────────┘
             │
             ├─► Jika ada baris terpengaruh (RowsAffected > 0):
             │   Tulis log: "Automatically expired N stale HOLD booking(s)"
             │
             ▼
     Tunggu Detak Ticker Berikutnya / Menunggu stopChan (Graceful Shutdown)
```

---

## 3. Implementasi Kode & Efisiensi Database

### Kueri Atomik Single-Statement
Alih-alih mengambil (*SELECT*) seluruh record ke memori lalu melakukan loop update satu per satu, worker mengeksekusi satu query SQL update langsung yang sangat cepat dan efisien:

```go
func ExpireStaleHoldBookings(db *gorm.DB) (int64, error) {
    now := time.Now()
    res := db.Model(&model.Booking{}).
        Where("status = ? AND hold_expires_at IS NOT NULL AND hold_expires_at < ?", model.StatusHold, now).
        Update("status", model.StatusCancelled)

    if res.Error != nil {
        return 0, res.Error
    }
    return res.RowsAffected, nil
}
```

### Penanganan Graceful Shutdown
Worker tidak akan menggantung atau mengalami kebocoran goroutine saat server dimatikan, karena dikendalikan oleh channel `stopChan`:

```go
func StartHoldExpiryWorker(db *gorm.DB, interval time.Duration) chan struct{} {
    stopChan := make(chan struct{})

    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                count, err := ExpireStaleHoldBookings(db)
                if err == nil && count > 0 {
                    log.Printf("[ExpiryWorker] Automatically expired %d stale HOLD booking(s)", count)
                }
            case <-stopChan:
                log.Println("[ExpiryWorker] Stopped background worker.")
                return
            }
        }
    }()

    return stopChan
}
```

Pada `cmd/server/main.go`, channel ini ditutup dengan elegan saat server berhenti:
```go
stopWorker := service.StartHoldExpiryWorker(database, expiryInterval)
defer close(stopWorker)
```
