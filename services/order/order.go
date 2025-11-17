package services

import (
	"context"
	"order-service/clients"
	"order-service/common/util"
	"order-service/domain/dto"
	"order-service/repositories"
)

type OrderService struct {
	repositories repositories.IRegistryRepository
	client       clients.IRegistryClient
}

type IOrderService interface {
	GetAllWithPagination(context.Context, *dto.OrderRequestParam) (*util.PaginationResult, error)
	GetByOrderID(context.Context, string) (*dto.OrderResponse, error)
	GetOrderByUserID(context.Context) ([]dto.OrderByUserIDResponse, error)
	Create(context.Context, *dto.OrderRequest) (*dto.OrderResponse, error)
	HandlePayment(context.Context, *dto.PaymentData) error
}

func NewOrderService(repositories repositories.IRegistryRepository, client clients.IRegistryClient) IOrderService {
	return &OrderService{
		repositories: repositories,
		client:       client,
	}
}

func (os *OrderService) GetAllWithPagination(context.Context, *dto.OrderRequestParam) (*util.PaginationResult, error) {
	panic("implement me")
}

func (os *OrderService) GetByOrderID(context.Context, string) (*dto.OrderResponse, error) {
	panic("implement me")
}

func (os *OrderService) GetOrderByUserID(context.Context) ([]dto.OrderByUserIDResponse, error) {
	panic("implement me")
}

func (os *OrderService) Create(context.Context, *dto.OrderRequest) (*dto.OrderResponse, error) {
	panic("implement me")
}

func (os *OrderService) HandlePayment(context.Context, *dto.PaymentData) error {
	panic("implement me")
}
