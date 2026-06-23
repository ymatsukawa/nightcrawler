package main

import (
	"time"
)

type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Email     string
	CreatedAt time.Time
	Orders    []Order `gorm:"foreignKey:UserID"`
}

type FilterUser struct {
	ID    uint `gorm:"primaryKey"`
	Email string
}

type Order struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Amount    int
	CreatedAt time.Time
	Items     []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID      uint `gorm:"primaryKey"`
	OrderID uint
	Product string
	Qty     int
}
