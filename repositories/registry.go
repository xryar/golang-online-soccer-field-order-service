package repositories

import (
	orderRepository "order-service/repositories/order"
	orderFieldRepository "order-service/repositories/orderField"
	orderHistoryRepository "order-service/repositories/orderHistory"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRegistryRepository interface {
	GetOrder() orderRepository.IOrderRepository
	GetOrderField() orderFieldRepository.IOrderFieldRepository
	GetOrderHistory() orderHistoryRepository.IOrderHistoryRepository
}

func NewRegistry(db *gorm.DB) IRegistryRepository {
	return &Registry{db: db}
}

func (r *Registry) GetOrder() orderRepository.IOrderRepository {
	return orderRepository.NewOrderRepository(r.db)
}

func (r *Registry) GetOrderField() orderFieldRepository.IOrderFieldRepository {
	return orderFieldRepository.NewOrderFieldRepository(r.db)
}

func (r *Registry) GetOrderHistory() orderHistoryRepository.IOrderHistoryRepository {
	return orderHistoryRepository.NewOrderHistoryRepository(r.db)
}
