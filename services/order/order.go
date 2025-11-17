package services

import (
	"context"
	"order-service/clients"
	"order-service/common/util"
	"order-service/domain/dto"
	"order-service/repositories"
)

type OrderService struct {
	repository repositories.IRegistryRepository
	client     clients.IRegistryClient
}

type IOrderService interface {
	GetAllWithPagination(context.Context, *dto.OrderRequestParam) (*util.PaginationResult, error)
	GetByOrderID(context.Context, string) (*dto.OrderResponse, error)
	GetOrderByUserID(context.Context) ([]dto.OrderByUserIDResponse, error)
	Create(context.Context, *dto.OrderRequest) (*dto.OrderResponse, error)
	HandlePayment(context.Context, *dto.PaymentData) error
}

func NewOrderService(repository repositories.IRegistryRepository, client clients.IRegistryClient) IOrderService {
	return &OrderService{
		repository: repository,
		client:     client,
	}
}

func (os *OrderService) GetAllWithPagination(ctx context.Context, param *dto.OrderRequestParam) (*util.PaginationResult, error) {
	orders, total, err := os.repository.GetOrder().FindAllWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}

	orderResult := make([]dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		user, err := os.client.GetUser().GetUserByUUID(ctx, order.UserId)
		if err != nil {
			return nil, err
		}

		orderResult = append(orderResult, dto.OrderResponse{
			UUID:      order.UUID,
			Code:      order.Code,
			Username:  user.Name,
			Amount:    order.Amount,
			Status:    order.Status.GetStatusString(),
			OrderDate: order.Date,
			CreatedAt: *order.CreatedAt,
			UpdatedAt: *order.UpdatedAt,
		})
	}

	paginationParam := util.PaginationParam{
		Page:  param.Page,
		Limit: param.Limit,
		Count: total,
		Data:  orderResult,
	}

	response := util.GeneratePagination(paginationParam)

	return &response, nil
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
