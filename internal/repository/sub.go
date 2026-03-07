package repository

import (
	"time"

	"github.com/google/uuid"
)

type Sub struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"size:255;not null" json:"service_name"`
	Price     int        `gorm:"not null" json:"price"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	StartDate time.Time  `gorm:"not null" json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}
