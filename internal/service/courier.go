package service

import (
	"context"
	"time"

	"github.com/doter2410/avito-project/internal/model"
)

// CourierRepository описывает, что мы ожидаем от базы данных
type CourierRepository interface {
	CreateCourier(ctx context.Context, c model.Courier) (int64, error)
	GetCourierById(ctx context.Context, id int64) (*model.Courier, error)
	GetCouriers(ctx context.Context) ([]*model.Courier, error)
	UpdateCourier(ctx context.Context, id int64, c model.Courier) error

	AssignCourierToOrder(ctx context.Context, orderID string, calcDeadline func(transportType string) time.Time) (*model.Courier, *model.Delivery, error)

	UnassignCourierFromOrder(ctx context.Context, orderID string) (*model.Courier, error)
}

type CourierService struct {
	repo    CourierRepository // Сервис зависит от интерфейса, а не от конкретной БД
	factory DeliveryTimeFactory
}

func NewCourierService(repo CourierRepository, factory DeliveryTimeFactory) *CourierService {
	return &CourierService{
		repo:    repo,
		factory: factory,
	}
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

func (s *CourierService) AssignDelivery(ctx context.Context, orderID string) (*model.Courier, *model.Delivery, error) {
	// Вызываем базу и передаем ей метод из фабрики в качестве аргумента
	return s.repo.AssignCourierToOrder(ctx, orderID, s.factory.CalculateDeadline)
}

func (s *CourierService) UnassignDelivery(ctx context.Context, orderID string) (*model.Courier, error) {
	return s.repo.UnassignCourierFromOrder(ctx, orderID)
}
