package repository

import (
	"log"
	"sync"
	"time"

	"dev.gaijin.team/go/golib/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubRepo - репозиторий для работы с подписками
type SubRepo struct {
	DB  *gorm.DB
	Log *logger.Logger
	mu  sync.Mutex
}

// NewSubRepo - создаёт новую запись сервиса
func NewSubRepo(db *gorm.DB) *SubRepo {
	return &SubRepo{
		DB: db,
		mu: sync.Mutex{},
	}
}

// AutoMigrate - авто миграции от gorm
func (r *SubRepo) AutoMigrate() error {
	log.Println("Запуск миграции...")

	if err := r.DB.Exec(`Создаем расширение если не существует "uuid-ossp"`).Error; err != nil {
		return err
	}

	err := r.DB.AutoMigrate(&Sub{})
	if err != nil {
		return err
	}

	log.Println("Миграция завершена!")
	return nil
}

func (r *SubRepo) Create(db *gorm.DB, sub Sub) *Sub {
	r.mu.Lock()
	defer r.mu.Unlock()

	db.Create(sub)

	return &sub
}

func (r *SubRepo) Reed(db *gorm.DB, id int64) *Sub {
	r.mu.Lock()
	defer r.mu.Unlock()

	var sub Sub

	db.First(&sub, id)

	return &sub
}

func (r *SubRepo) Update(db *gorm.DB, sub Sub) *Sub {
	r.mu.Lock()
	defer r.mu.Unlock()

	db.Save(&sub)

	return &sub
}

func (r *SubRepo) Delete(db *gorm.DB, id int64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var sub Sub
	db.First(&sub, id)

	db.Delete(&sub)

}

func (r *SubRepo) GetPriceForRange(db *gorm.DB, idUser uuid.UUID, idSub uuid.UUID, startData time.Time, endData time.Time) int64 {
	r.mu.Lock()
	defer r.mu.Unlock()

	var prices int64
	db.Model(&Sub{}).
		Where(&Sub{ID: idSub, UserID: idUser}).
		Where("start_date >= ?", startData).
		Where("end_date <= ?", endData).
		Select("COALESCE(SUM(price), 0)").
		Row().
		Scan(&prices)

	return prices
}
