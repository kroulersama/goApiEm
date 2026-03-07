package service

import (
	"goApiEM/internal/repository"
	"time"

	"github.com/google/uuid"
)

// SubRepo - интерфейс для работы с репозиторием подписок
type SubRepo interface {
	Create(sub repository.Sub) *repository.Sub
	GetByID(id int64) *repository.Sub
	Update(sub repository.Sub) *repository.Sub
	Delete(id int64)
	GetPriceForRange(idSub int64, idUser uuid.UUID,
		startData time.Time, endData time.Time) int64
}

// subTestRepo - тестовый репозиторий
type subTestRepo struct {
}

func (r *subTestRepo) Create(sub repository.Sub) *repository.Sub {
	return &repository.Sub{}
}

func (r *subTestRepo) GetByID(id int64) *repository.Sub {
	return &repository.Sub{}
}
func (r *subTestRepo) Update(sub repository.Sub) *repository.Sub {
	return &repository.Sub{}
}
func (r *subTestRepo) Delete(id int64) {}

func (r *subTestRepo) GetPriceForRange(idSub int64, idUser uuid.UUID,
	startData time.Time, endData time.Time) int64 {
	return 0
}

// SubService - бизнесс-логика для работы с подписками
type SubService struct {
	repo SubRepo
}

// NewSubSevrice - Создание нового сервиса
func NewSubSevrice(repo SubRepo) *SubService { return &SubService{repo: repo} }

// CreateSub - создание новой подписки
func (s *SubService) CreateSub(name string, price int, userID uuid.UUID,
	startDate time.Time, endDate *time.Time) (repository.Sub, error) {

	sub := repository.Sub{
		Name:      name,
		Price:     price,
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	s.repo.Create(sub)

	return sub, nil

}

// ReedByIDSub - чтение записи подписки по id
func (s *SubService) GetByIDSub(id int64) (*repository.Sub, error) {

	sub := s.repo.GetByID(id)

	return sub, nil
}

// UpdateSub - обновляет запись подписки
func (s *SubService) UpdateSub(id int64, name string, price int, userID uuid.UUID,
	startDate time.Time, endDate *time.Time) (repository.Sub, error) {

	sub := repository.Sub{
		ID:        id,
		Name:      name,
		Price:     price,
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}
	s.repo.Update(sub)

	return sub, nil
}

// DeleteSub - удаление записи подписки по id
func (s *SubService) DeleteSub(id int64) error {

	s.repo.Delete(id)

	return nil
}

// GetPriceForRangeSub - подсчёт суммы подписок за период времени по id подписки и id юзера
func (s *SubService) GetPriceForRangeSub(idSub int64, idUser uuid.UUID,
	startData time.Time, endData time.Time) (int64, error) {

	prices := s.repo.GetPriceForRange(idSub, idUser, startData, endData)

	return prices, nil
}
