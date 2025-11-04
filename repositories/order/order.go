package repositories

import (
	"context"
	"order-service/domain/dto"
	"order-service/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

type IOrderRepository interface {
	FindAllWithPagination(context.Context, *dto.OrderRequestParam) ([]models.Order, int64, error)
	FindByUserID(context.Context, string) ([]models.Order, error)
	Create(context.Context, *gorm.DB, *models.Order) (*models.Order, error)
	Update(context.Context, *gorm.DB, string, *models.Order, uuid.UUID) error
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{db: db}
}

func (or *OrderRepository) FindAllWithPagination(context.Context, *dto.OrderRequestParam) ([]models.Order, int64, error) {
}

func (or *OrderRepository) FindByUserID(context.Context, string) ([]models.Order, error) {
}

func (or *OrderRepository) Create(context.Context, *gorm.DB, *models.Order) (*models.Order, error) {
}

func (or *OrderRepository) Update(context.Context, *gorm.DB, string, *models.Order, uuid.UUID) error {
}
