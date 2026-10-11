package service

import (
	"log"
	"time"

	"github.com/raygbrn/ai-booking-hotel-system/internal/model"
	"gorm.io/gorm"
)

// ExpireStaleHoldBookings finds all bookings with status HOLD whose hold_expires_at is in the past
// and cancels them, releasing room units back into inventory.
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

// StartHoldExpiryWorker runs a background ticker that periodically cancels expired HOLD reservations.
// Returns a stop channel to terminate gracefully.
func StartHoldExpiryWorker(db *gorm.DB, interval time.Duration) chan struct{} {
	stopChan := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				count, err := ExpireStaleHoldBookings(db)
				if err != nil {
					log.Printf("[ExpiryWorker] Error expiring stale bookings: %v", err)
				} else if count > 0 {
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
