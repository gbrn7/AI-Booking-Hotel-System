# Modul 04: Engine Inventori Kamar & Proteksi Overbooking Atomik

## 1. Tantangan Bisnis Inventori Perhotelan

Berbeda dengan sistem reservasi klinik (seperti referensi `AI-Sales-Agent-Workshop-Go`) yang berbasis **slot jam** (misal slot 90 menit pada hari tertentu), sistem hotel beroperasi pada prinsip **rentang malam menginap (Date-Range Capacity)**.

Tamu memesan kamar dari tanggal Check-In ($D_{in}$) hingga Check-Out ($D_{out}$), yang berarti tamu menempati kamar pada setiap malam di dalam rentang:
$$\{D_{in}, D_{in}+1, \dots, D_{out}-1\}$$

Tantangan terbesarnya adalah: sebuah kamar hanya dapat dinyatakan tersedia jika pada **seluruh malam** dalam rentang menginap tersebut, sisa unit kamar $\ge$ jumlah kamar yang diminta tamu.

---

## 2. Formula Kalkulasi Ketersediaan

Pada fungsi `calculateAvailableUnitsForRoom`:
1. Hitung total malam menginap: $\text{totalNights} = \text{DateDiff}(D_{out}, D_{in})$.
2. Untuk setiap malam $N \in [0, \text{totalNights}-1]$:
   - Tanggal malam ke-$N$: $D_N = D_{in} + N \text{ hari}$.
   - Kueri seluruh booking aktif yang menempati malam $D_N$:
   ```sql
   SELECT COALESCE(SUM(total_rooms), 0)
   FROM bookings
   WHERE room_id = ?
     AND status IN ('HOLD', 'CONFIRMED', 'CHECKED_IN')
     AND check_in_date <= D_N
     AND check_out_date > D_N;
   ```
3. Cari titik puncak keterisian tertinggi:
   $$\text{maxBooked} = \max_{N} (\text{bookedUnits}(D_N))$$
4. Sisa unit yang aman untuk disewa:
   $$\text{availableUnits} = \max(0, \text{totalUnits} - \text{maxBooked})$$

---

## 3. Pencegahan Overbooking Simultan (Concurrency Protection)

### Skenario Ancaman (Race Condition)
Bayangkan kamar tipe *Executive Ocean Suite* hanya tersisa **1 unit** terakhir untuk akhir pekan (15–17 Oktober).
* Tamu A meminta booking lewat WhatsApp di detik `10:00:01.100`.
* Tamu B meminta booking lewat Web Widget di detik `10:00:01.120`.

Tanpa sinkronisasi, kedua thread membaca sisa unit = 1, keduanya membuat status `HOLD`, sehingga terjadi **double-booking (overbooking)** yang merugikan reputasi hotel.

### Solusi: Mutex Registry & Database Transaction
Engine inventori di `internal/service/inventory.go` mengimplementasikan strategi perlindungan berlapis:

```go
// 1. Mutex Registry Per-Room ID
var (
    roomMutexMap   = make(map[uint]*sync.Mutex)
    roomRegistryMu sync.Mutex
)

func getRoomMutex(roomID uint) *sync.Mutex {
    roomRegistryMu.Lock()
    defer roomRegistryMu.Unlock()
    if mu, exists := roomMutexMap[roomID]; exists {
        return mu
    }
    mu := &sync.Mutex{}
    roomMutexMap[roomID] = mu
    return mu
}
```

Saat `HoldBooking` dipanggil:
1. Backend mengunci thread mutex untuk kamar tersebut (`mu.Lock()`).
2. Seluruh evaluasi ketersediaan dan pembuatan record booking dibungkus dalam transaksi database atomik:
   ```go
   err = db.Transaction(func(tx *gorm.DB) error {
       // Kueri ulang ketersediaan di dalam transaksi
       availableUnits, err := calculateAvailableUnitsForRoom(tx, room.ID, room.TotalUnits, checkIn, totalNights)
       if availableUnits < req.TotalRooms {
           return ErrNoRoomsAvailable
       }
       // Simpan booking status HOLD
       return tx.Create(&newBooking).Error
   })
   ```
3. Thread melepaskan mutex (`defer mu.Unlock()`).

### Hasil Uji Konkurensi
Unit test `TestHoldBookingConcurrency` di `internal/service/inventory_test.go` membuktikan bahwa ketika 5 goroutine mencoba memesan 1 unit kamar terakhir secara bersamaan, **tepat 1 reservasi yang berhasil**, dan 4 lainnya menerima error `ErrNoRoomsAvailable`.
