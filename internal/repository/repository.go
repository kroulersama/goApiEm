package repository

import (
	"goApiEM/internal/model"
	"log"
	"sync"

	"dev.gaijin.team/go/golib/logger"
	"gorm.io/gorm"
)

// SubRepo - репозитоия для работы с подписками
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

func (r *SubRepo) AutoMigrate() error {
	log.Println("Запуск миграции...")

	if err := r.DB.Exec(`Создаем расширение если не существует "uuid-ossp"`).Error; err != nil {
		return err
	}

	err := r.DB.AutoMigrate(&model.Sub{})
	if err != nil {
		return err
	}

	log.Println("Миграция завершена!")
	return nil
}

func Create() {
	// todo
}

func Reed() {
	// todo
}

func Update() {
	// todo
}

func Delete() {
	//todo
}

func GetAllPrice() {
	//todo
}
