package service

import (
	"goApiEM/internal/repository"
	"time"

	"github.com/google/uuid"
)

// SubRepo - интерфейс для работы с репозиторием подписок
type SubRepo interface {
	Create(sub repository.Sub) *repository.Sub
	ReedByID(id int64) *repository.Sub
	Update(sub repository.Sub) *repository.Sub
	Delete(id int64)
	GetPriceForRange(idUser uuid.UUID, idSub uuid.UUID,
		startData time.Time, endData time.Time) int64
}

// SubService - бизнесс-логика для работы с подписками
type SubService struct {
	repo SubRepo
}

// NewSubSevrice - Создание нового сервиса
func NewSubSevrice(repo SubRepo) *SubService { return &SubService{repo: repo} }

// CreateSub - создание новой подписки
func (s *SubService) CreateSub(name string, price int, userID uuid.UUID,
	startDate time.Time, endDate *time.Time) repository.Sub {

	sub := repository.Sub{
		Name:      name,
		Price:     price,
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	s.repo.Create(sub)

	return sub

}

// ReedByIDSub - чтение записи подписки по id
func (s *SubService) ReedByIDSub(id int64) *repository.Sub {

	sub := s.repo.ReedByID(id)

	return sub
}

// UpdateSub - обновляет запись подписки
func (s *SubService) UpdateSub(id uuid.UUID, name string, price int, userID uuid.UUID,
	startDate time.Time, endDate *time.Time) repository.Sub {

	sub := repository.Sub{
		ID:        id,
		Name:      name,
		Price:     price,
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}
	s.repo.Update(sub)

	return sub
}

// DeleteSub - удаление записи подписки по id
func (s *SubService) DeleteSub(id int64) int64 {

	s.repo.Delete(id)

	return id
}

// GetPriceForRangeSub - подсчёт суммы подписок за период времени по id подписки и id юзера
func (s *SubService) GetPriceForRangeSub(idUser uuid.UUID, idSub uuid.UUID,
	startData time.Time, endData time.Time) int64 {

	prices := s.repo.GetPriceForRange(idUser, idSub, startData, endData)

	return prices
}
