package repository

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("запись не найдена")
var ErrInvalidId = errors.New("id повторился")
var ErrReqDB = errors.New("ошибка обращения к базе данных")

// SubRepo - репозиторий для работы с подписками
type SubRepo struct {
	DB *gorm.DB
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
func (r *SubRepo) Create(sub Sub) error {

	if err := r.DB.Create(&sub).Error; err != nil {
		return ErrInvalidId
	}

	return nil
}

// ReedById - чтение одной записи по id
func (r *SubRepo) GetByID(id int64) (*Sub, error) {

	var sub Sub

	if err := r.DB.First(&sub, id).Error; err != nil {
		return nil, ErrNotFound
	}

	return &sub, nil
}

// Update - обновление инфолмации записи
func (r *SubRepo) Update(sub Sub) error {

	if err := r.DB.Save(&sub).Error; err != nil {
		return ErrNotFound
	}

	return nil
}

// Delete - удалиние записи из базы по id
func (r *SubRepo) Delete(id int64) error {

	var sub Sub
	if err := r.DB.First(&sub, id).Error; err != nil {
		return ErrNotFound
	}

	if err := r.DB.Delete(&sub).Error; err != nil {
		return ErrReqDB
	}

	return nil

}

// List - возвращает все подписки
func (r *SubRepo) List() (*[]Sub, error) {

	var subs []Sub

	if err := r.DB.Find(&subs).Error; err != nil {
		return nil, ErrReqDB
	}

	return &subs, nil
}

// GetPriceForRange - выявление суммы подписки за период времени по id подписки, и id юзера
func (r *SubRepo) GetPriceForRange(idSub int64, idUser uuid.UUID,
	startData time.Time, endData time.Time) (int64, error) {

	var prices int64

	var subName string
	err := r.DB.Model(&Sub{}).
		Where("id = ?", idSub).
		Select("name").
		Scan(&subName).Error
	if err != nil {
		return 0, ErrReqDB
	}

	err = r.DB.Model(&Sub{}).
		Where("name = ?", subName).
		Where("start_date >= ?", startData).
		Where("end_date <= ?", endData).
		Select("COALESCE(SUM(price), 0)").
		Row().
		Scan(&prices)
	if err != nil {
		return 0, ErrReqDB
	}

	return prices, nil
}
