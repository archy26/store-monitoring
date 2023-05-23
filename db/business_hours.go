package db

type StoreBusinessHours struct {
	StoreID        string `gorm:"primaryKey"`
	Day            int    `gorm:"primaryKey"`
	StartTimeLocal string `gorm:"primaryKey"`
	EndTimeLocal   string `gorm:"not null"`
}
