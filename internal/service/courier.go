package service

import (
	"context"

	"github.com/doter2410/avito-project/internal/model"
)

// CourierRepository описывает, что мы ожидаем от базы данных
type CourierRepository interface {
	CreateCourier(ctx context.Context, c model.Courier) (int64, error)
	GetCourierById(ctx context.Context, id int64) (*model.Courier, error)
	GetCouriers(ctx context.Context) ([]*model.Courier, error)
	UpdateCourier(ctx context.Context, id int64, c model.Courier) error
}

type CourierService struct {
	repo CourierRepository // Сервис зависит от интерфейса, а не от конкретной БД
}

func NewCourierService(repo CourierRepository) *CourierService {
	return &CourierService{repo: repo}
}

// Ниже идут методы сервиса. Пока они просто вызывают репозиторий,
// но в будущем здесь будет логика (например, проверка статусов)

func (s *CourierService) CreateCourier(ctx context.Context, c model.Courier) (int64, error) {
	return s.repo.CreateCourier(ctx, c)
}

func (s *CourierService) GetCourierById(ctx context.Context, id int64) (*model.Courier, error) {
	return s.repo.GetCourierById(ctx, id)
}

func (s *CourierService) GetCouriers(ctx context.Context) ([]*model.Courier, error) {
	return s.repo.GetCouriers(ctx)
}

func (s *CourierService) UpdateCourier(ctx context.Context, id int64, c model.Courier) error {
	return s.repo.UpdateCourier(ctx, id, c)
}
