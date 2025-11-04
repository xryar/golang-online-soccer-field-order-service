package repositories

import (
	"context"
	errWrap "order-service/common/error"
	errConstant "order-service/constants/error"
	"order-service/domain/dto"
	"order-service/domain/models"

	"gorm.io/gorm"
)

type OrderHistoryRepository struct {
	db *gorm.DB
}

type IOrderHistoryRepository interface {
	Create(context.Context, *gorm.DB, *dto.OrderHistoryRequest) error
}

func NewOrderHistoryRepository(db *gorm.DB) IOrderHistoryRepository {
	return &OrderHistoryRepository{db: db}
}

func (or *OrderHistoryRepository) Create(ctx context.Context, tx *gorm.DB, req *dto.OrderHistoryRequest) error {
	orderHistory := &models.OrderHistory{
		OrderID: req.OrderID,
		Status:  req.Status,
	}

	err := tx.WithContext(ctx).Create(&orderHistory).Error
	if err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	return nil
}
