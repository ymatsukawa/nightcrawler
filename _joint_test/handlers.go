package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func selectManyHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var users []User
		if err := db.Find(&users).Error; err != nil {
			return err
		}
		return c.JSON(http.StatusOK, echo.Map{"category": "Select many", "users": len(users)})
	}
}

func indexNoWorthHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var users []User
		if err := db.Select("id, name").Where("name LIKE ?", "%a").Limit(10).Find(&users).Error; err != nil {
			return err
		}
		return c.JSON(http.StatusOK, echo.Map{"category": "Index no worth", "users": len(users)})
	}
}

func subqueryHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sub := db.Model(&Order{}).Select("user_id")
		var users []User
		if err := db.Select("id, name").Where("id IN (?)", sub).Limit(10).Find(&users).Error; err != nil {
			return err
		}
		return c.JSON(http.StatusOK, echo.Map{"category": "Using subquery", "users": len(users)})
	}
}

func heavyCalcHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var rows []map[string]interface{}
		if err := db.Table("users").
			Select("users.id").
			Joins("JOIN orders ON orders.user_id = users.id").
			Joins("JOIN order_items ON order_items.order_id = orders.id").
			Limit(10).
			Scan(&rows).Error; err != nil {
			return err
		}
		return c.JSON(http.StatusOK, echo.Map{"category": "Heavy clause", "rows": len(rows)})
	}
}

func usingFunctionHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var users []User
		if err := db.Select("id, name").Where("DATE(created_at) = ?", "2024-01-01").Limit(10).Find(&users).Error; err != nil {
			return err
		}
		return c.JSON(http.StatusOK, echo.Map{"category": "Using function", "users": len(users)})
	}
}

func nPlusOneHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var users []User
		if err := db.Find(&users).Error; err != nil {
			return err
		}
		total := 0
		for _, u := range users {
			var orders []Order
			if err := db.Where("user_id = ?", u.ID).Find(&orders).Error; err != nil {
				return err
			}
			total += len(orders)
		}
		return c.JSON(http.StatusOK, echo.Map{"category": "N+1", "users": len(users), "orders": total})
	}
}

func noSloQueryHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user []FilterUser
		if err := db.Model(&User{}).Where("id = ?", 1).Limit(1).Find(&user).Error; err != nil {
			return err
		}
		return c.JSON(http.StatusOK, echo.Map{"category": "No slow query", "user": len(user)})
	}
}
