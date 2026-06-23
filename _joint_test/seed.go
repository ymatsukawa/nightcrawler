package main

import (
	"fmt"

	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func migrateAndSeed(db *gorm.DB) error {
	silent := db.Session(&gorm.Session{Logger: glogger.Discard})

	if err := silent.AutoMigrate(&User{}, &Order{}, &OrderItem{}); err != nil {
		return fmt.Errorf("automigrate: %w", err)
	}

	var count int64
	silent.Model(&User{}).Count(&count)
	if count > 0 {
		return nil
	}

	for u := 1; u <= 5; u++ {
		user := User{
			Name:  fmt.Sprintf("user_%d", u),
			Email: fmt.Sprintf("user%d@example.com", u),
		}
		if err := silent.Create(&user).Error; err != nil {
			return fmt.Errorf("seed user: %w", err)
		}
		for o := 1; o <= 3; o++ {
			order := Order{UserID: user.ID, Amount: u * o * 100}
			if err := silent.Create(&order).Error; err != nil {
				return fmt.Errorf("seed order: %w", err)
			}
			for i := 1; i <= 2; i++ {
				item := OrderItem{
					OrderID: order.ID,
					Product: fmt.Sprintf("product_%d_%d", o, i),
					Qty:     i,
				}
				if err := silent.Create(&item).Error; err != nil {
					return fmt.Errorf("seed item: %w", err)
				}
			}
		}
	}
	return nil
}
