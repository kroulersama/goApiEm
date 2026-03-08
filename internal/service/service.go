package service

import (
	"goApiEM/internal/repository"
	"time"

	"github.com/google/uuid"
)

// SubRepo - интерфейс для работы с репозиторием подписок
type SubRepo interface {
	Create(sub repository.Sub) error
	GetByID(id int64) (*repository.Sub, error)
	Update(sub repository.Sub) error
	Delete(id int64) error
	List() (*[]repository.Sub, error)
	GetPriceForRange(idSub int64, idUser uuid.UUID,
		startData time.Time, endData time.Time) (int64, error)
}

// subTestRepo - тестовый репозиторий
type subTestRepo struct{}

func (r *subTestRepo) Create(sub repository.Sub) error           { return nil }
func (r *subTestRepo) GetByID(id int64) (*repository.Sub, error) { return &repository.Sub{}, nil }
func (r *subTestRepo) Update(sub repository.Sub) error           { return nil }
func (r *subTestRepo) Delete(id int64) error                     { return nil }
func (r *subTestRepo) List() (*[]repository.Sub, error)          { return &[]repository.Sub{}, nil }
func (r *subTestRepo) GetPriceForRange(idSub int64, idUser uuid.UUID,
	startData time.Time, endData time.Time) (int64, error) {
	return 0, nil
}

// SubService - бизнесс-логика для работы с подписками
type SubService struct {
	repo SubRepo
}

// NewSubSevrice - Создание нового сервиса
func NewSubSevrice(repo SubRepo) *SubService { return &SubService{repo: repo} }

// CreateSub - создание новой подписки
func (s *SubService) CreateSub(name string, price int, userID uuid.UUID,
	startDate time.Time, endDate *time.Time) (*repository.Sub, error) {

	sub := repository.Sub{
		Name:      name,
		Price:     price,
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	if err := s.repo.Create(sub); err != nil {
		return nil, err
	}

	return &sub, nil

}

// ReedByIDSub - чтение записи подписки по id
func (s *SubService) GetByIDSub(id int64) (*repository.Sub, error) {

	sub, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

// UpdateSub - обновляет запись подписки
func (s *SubService) UpdateSub(id int64, name string, price int, userID uuid.UUID,
	startDate time.Time, endDate *time.Time) (*repository.Sub, error) {

	sub := repository.Sub{
		ID:        id,
		Name:      name,
		Price:     price,
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	if err := s.repo.Update(sub); err != nil {
		return nil, err
	}

	return &sub, nil
}

// DeleteSub - удаление записи подписки по id
func (s *SubService) DeleteSub(id int64) error {

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	return nil
}

// GetSubs - возвращает все подписки
func (s *SubService) GetSubs() (*[]repository.Sub, error) {

	subs, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	return subs, nil
}

// GetPriceForRangeSub - подсчёт суммы подписок за период времени по id подписки и id юзера
func (s *SubService) GetPriceForRangeSub(idSub int64, idUser uuid.UUID,
	startData time.Time, endData time.Time) (int64, error) {

	prices, err := s.repo.GetPriceForRange(idSub, idUser, startData, endData)
	if err != nil {
		return 0, err
	}

	return prices, nil
}
