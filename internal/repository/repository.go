package repository

import (
	"log"
	"time"

	"dev.gaijin.team/go/golib/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubRepo - репозиторий для работы с подписками
type SubRepo struct {
	DB  *gorm.DB
	Log *logger.Logger
}

// NewSubRepo - создаёт новую запись сервиса
func NewSubRepo(db *gorm.DB) *SubRepo {
	return &SubRepo{
		DB: db,
	}
}

// AutoMigrate - авто миграции от gorm
func (r *SubRepo) AutoMigrate() error {
	log.Println("Запуск миграции...")

	if err := r.DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return err
	}

	err := r.DB.AutoMigrate(&Sub{})
	if err != nil {
		return err
	}

	log.Println("Миграция завершена!")
	return nil
}

// Create - создание записи в базе
func (r *SubRepo) Create(sub Sub) *Sub {

	r.DB.Create(sub)

	return &sub
}

// ReedById - чтение одной записи по id
func (r *SubRepo) ReedByID(id int64) *Sub {

	var sub Sub

	r.DB.First(&sub, id)

	return &sub
}

// Update - обновление инфолмации записи
func (r *SubRepo) Update(sub Sub) *Sub {

	r.DB.Save(&sub)

	return &sub
}

// Delete - удалиние записи из базы по id
func (r *SubRepo) Delete(id int64) {

	var sub Sub
	r.DB.First(&sub, id)

	r.DB.Delete(&sub)

}

// GetPriceForRange - выявление суммы подписки за период времени по id подписки, и id юзера
func (r *SubRepo) GetPriceForRange(idUser uuid.UUID, idSub uuid.UUID,
	startData time.Time, endData time.Time) int64 {

	var prices int64
	r.DB.Model(&Sub{}).
		Where(&Sub{ID: idSub, UserID: idUser}).
		Where("start_date >= ?", startData).
		Where("end_date <= ?", endData).
		Select("COALESCE(SUM(price), 0)").
		Row().
		Scan(&prices)

	return prices
}
