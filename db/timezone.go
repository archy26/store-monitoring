package db

type Timezone struct {
	ID          int    `gorm:"primaryKey;autoIncrement"`
	StoreID     string `gorm:"not null"`
	TimezoneStr string `gorm:"not null"`
}
