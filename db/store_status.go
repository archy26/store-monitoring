package db

import "time"

type StoreData struct {
	StoreID      string `gorm:"primaryKey"`
	Status       string
	TimestampUTC time.Time `gorm:"primaryKey"`
}
