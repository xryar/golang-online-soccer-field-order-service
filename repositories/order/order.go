package repositories

import (
	"context"
	"fmt"
	errWrap "order-service/common/error"
	errConstant "order-service/constants/error"
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
	FindByUUID(context.Context, string) (*models.Order, error)
	Create(context.Context, *gorm.DB, *models.Order) (*models.Order, error)
	Update(context.Context, *gorm.DB, string, *models.Order, uuid.UUID) error
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{db: db}
}

func (or *OrderRepository) FindAllWithPagination(ctx context.Context, param *dto.OrderRequestParam) ([]models.Order, int64, error) {
	var (
		orders []models.Order
		sort   string
		total  int64
	)
	if param.SortColumn != nil {
		sort = fmt.Sprintf("%s %s", *param.SortColumn, *param.SortOrder)
	} else {
		sort = "created_at desc"
	}

	limit := param.Limit
	offset := (param.Page - 1) * limit
	err := or.db.WithContext(ctx).Limit(limit).Offset(offset).Order(sort).Find(&orders).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	err = or.db.WithContext(ctx).Model(&models.Order{}).Count(&total).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return orders, total, nil
}

func (or *OrderRepository) FindByUserID(context.Context, string) ([]models.Order, error) {
}

func (or *OrderRepository) FindByUUID(context.Context, string) (*models.Order, error) {

}

func (or *OrderRepository) Create(context.Context, *gorm.DB, *models.Order) (*models.Order, error) {
}

func (or *OrderRepository) Update(context.Context, *gorm.DB, string, *models.Order, uuid.UUID) error {
}
