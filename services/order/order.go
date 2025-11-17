package services

import (
	"context"
	"fmt"
	"order-service/clients"
	clientUser "order-service/clients/user"
	"order-service/common/util"
	"order-service/constants"
	"order-service/domain/dto"
	"order-service/domain/models"
	"order-service/repositories"
)

type OrderService struct {
	repository repositories.IRegistryRepository
	client     clients.IRegistryClient
}

type IOrderService interface {
	GetAllWithPagination(context.Context, *dto.OrderRequestParam) (*util.PaginationResult, error)
	GetByUUID(context.Context, string) (*dto.OrderResponse, error)
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

func (os *OrderService) GetByUUID(ctx context.Context, uuid string) (*dto.OrderResponse, error) {
	var (
		order *models.Order
		user  *clientUser.UserData
		err   error
	)

	order, err = os.repository.GetOrder().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	user, err = os.client.GetUser().GetUserByUUID(ctx, order.UserId)
	if err != nil {
		return nil, err
	}

	response := dto.OrderResponse{
		UUID:      order.UUID,
		Code:      order.Code,
		Username:  user.Name,
		Amount:    order.Amount,
		Status:    order.Status.GetStatusString(),
		OrderDate: order.Date,
		CreatedAt: *order.CreatedAt,
		UpdatedAt: *order.UpdatedAt,
	}

	return &response, nil
}

func (os *OrderService) GetOrderByUserID(ctx context.Context) ([]dto.OrderByUserIDResponse, error) {
	var (
		order []models.Order
		user  = ctx.Value(constants.User).(*clientUser.UserData)
		err   error
	)
	order, err = os.repository.GetOrder().FindByUserID(ctx, user.UUID.String())
	if err != nil {
		return nil, err
	}

	orderLists := make([]dto.OrderByUserIDResponse, 0, len(order))
	for _, item := range order {
		payment, err := os.client.GetPayment().GetPaymentByUUID(ctx, item.PaymentID)
		if err != nil {
			return nil, err
		}

		orderLists = append(orderLists, dto.OrderByUserIDResponse{
			Code:        item.Code,
			Amount:      fmt.Sprintf("%s", util.RupiahFormat(&item.Amount)),
			Status:      item.Status.GetStatusString(),
			OrderDate:   item.Date.String(),
			PaymentLink: payment.PaymentLink,
			InvoiceLink: payment.InvoiceLink,
		})
	}

	return orderLists, nil
}

func (os *OrderService) Create(context.Context, *dto.OrderRequest) (*dto.OrderResponse, error) {
	panic("implement me")
}

func (os *OrderService) HandlePayment(context.Context, *dto.PaymentData) error {
	panic("implement me")
}
