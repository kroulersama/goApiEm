package service

import (
	"gorm.io/gorm"
)

type SubRepo interface {
	Create(db *gorm.DB, sub Sub) *Sub
	Reed(db *gorm.DB, id int64) *Sub
	Update(db *gorm.DB, sub Sub) *Sub
	Delete(db *gorm.DB, id int64)
}

type SubService struct {
	repo SubRepo
}
