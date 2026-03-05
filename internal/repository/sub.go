package repository

import (
	"time"

	"github.com/google/uuid"
)

type Sub struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string     `gorm:"size:255;not null" json:"name"`
	Price     int        `gorm:"not null" json:"price"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	StartDate time.Time  `gorm:"not null" json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}
