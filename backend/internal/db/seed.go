package db

import (
	"log"

	"github.com/raygbrn/ai-booking-hotel-system/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedDefaultData inserts initial rooms and staff accounts if empty.
func SeedDefaultData(database *gorm.DB) error {
	var count int64
	database.Model(&model.Room{}).Count(&count)

	if count == 0 {
		rooms := []model.Room{
			{
				Name:              "Superior King Room",
				Type:              "standard",
				BasePricePerNight: 750000,
				CapacityGuests:    2,
				TotalUnits:        10,
				Description:       "Kamar nyaman dengan ranjang King Bed, smart TV 43 inci, pemandangan kota, dan kamar mandi modern.",
				Amenities:         `["King Bed", "Smart TV 43\"", "City View", "Free High-Speed WiFi", "Coffee/Tea Maker", "Rain Shower"]`,
				ImageURLs:         `["https://images.unsplash.com/photo-1590490360182-c33d57733427?auto=format&fit=crop&w=800&q=80"]`,
				IsActive:          true,
			},
			{
				Name:              "Deluxe Ocean View",
				Type:              "deluxe",
				BasePricePerNight: 1250000,
				CapacityGuests:    2,
				TotalUnits:        8,
				Description:       "Kamar elegan dengan balkon pribadi yang menghadap langsung ke lautan biru, bathtub marmer, dan sarapan buffet.",
				Amenities:         `["Private Balcony", "Direct Ocean View", "Marble Bathtub", "King Bed", "Buffet Breakfast", "Minibar"]`,
				ImageURLs:         `["https://images.unsplash.com/photo-1582719478250-c89cae4dc85b?auto=format&fit=crop&w=800&q=80"]`,
				IsActive:          true,
			},
			{
				Name:              "Executive Ocean Suite",
				Type:              "suite",
				BasePricePerNight: 2200000,
				CapacityGuests:    3,
				TotalUnits:        4,
				Description:       "Suite mewah dengan ruang tamu terpisah, jacuzzi pribadi di teras, akses exclusive lounge, dan pemandangan sunset spektakuler.",
				Amenities:         `["Private Jacuzzi", "Separate Living Room", "Executive Lounge Access", "Sunset View", "Complimentary High Tea", "King Bed"]`,
				ImageURLs:         `["https://images.unsplash.com/photo-1631049307264-da0ec9d70304?auto=format&fit=crop&w=800&q=80"]`,
				IsActive:          true,
			},
			{
				Name:              "Grand Family Villa",
				Type:              "family_suite",
				BasePricePerNight: 3500000,
				CapacityGuests:    5,
				TotalUnits:        2,
				Description:       "Villa keluarga dua lantai dengan kolam renang pribadi, dapur mini, 2 kamar tidur terpisah, dan ruang santai keluarga luas.",
				Amenities:         `["Private Plunge Pool", "2 Bedrooms (King + Twin)", "Kitchenette", "Garden & Pool View", "Dining Area", "Complimentary Airport Transfer"]`,
				ImageURLs:         `["https://images.unsplash.com/photo-1566073771259-6a8506099945?auto=format&fit=crop&w=800&q=80"]`,
				IsActive:          true,
			},
		}

		for _, r := range rooms {
			if err := database.Create(&r).Error; err != nil {
				return err
			}
		}
		log.Printf("Seeded %d hotel rooms successfully.", len(rooms))
	}

	var userCount int64
	database.Model(&model.User{}).Count(&userCount)
	if userCount == 0 {
		hashAdmin, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		hashStaff, _ := bcrypt.GenerateFromPassword([]byte("staff123"), bcrypt.DefaultCost)

		users := []model.User{
			{
				Username:     "admin",
				Email:        "admin@grandoasis.com",
				PasswordHash: string(hashAdmin),
				Name:         "General Manager",
				Role:         model.RoleAdmin,
			},
			{
				Username:     "receptionist",
				Email:        "receptionist@grandoasis.com",
				PasswordHash: string(hashStaff),
				Name:         "Front Desk Staff",
				Role:         model.RoleReceptionist,
			},
		}

		for _, u := range users {
			if err := database.Create(&u).Error; err != nil {
				return err
			}
		}
		log.Printf("Seeded %d staff users successfully.", len(users))
	}

	return nil
}
