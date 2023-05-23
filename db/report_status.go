package db

type ReportStatus struct {
	ReportID string `gorm:"primaryKey"`
	Status   string
}
