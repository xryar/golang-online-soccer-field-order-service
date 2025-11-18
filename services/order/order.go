package services

import (
	"context"
	"fmt"
	"order-service/clients"
	clientField "order-service/clients/field"
	clientPayment "order-service/clients/payment"
	clientUser "order-service/clients/user"
	"order-service/common/util"
	"order-service/constants"
	errOrder "order-service/constants/error/order"
	"order-service/domain/dto"
	"order-service/domain/models"
	"order-service/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (os *OrderService) Create(ctx context.Context, request *dto.OrderRequest) (*dto.OrderResponse, error) {
	var (
		order              *models.Order
		txErr, err         error
		user               = ctx.Value(constants.User).(*clientUser.UserData)
		field              *clientField.FieldData
		paymentResponse    *clientPayment.PaymentData
		orderFieldSchedule = make([]models.OrderField, 0, len(request.FieldScheduleIDs))
		totalAmount        float64
	)

	for _, fieldId := range request.FieldScheduleIDs {
		uuidParsed := uuid.MustParse(fieldId)
		field, err = os.client.GetField().GetFieldByUUID(ctx, uuidParsed)
		if err != nil {
			return nil, err
		}

		totalAmount += field.PricePerHour
		if field.Status == constants.BookedStatus.String() {
			return nil, errOrder.ErrFieldAlreadyBooked
		}
	}

	err = os.repository.GetTX().Transaction(func(tx *gorm.DB) error {
		order, txErr = os.repository.GetOrder().Create(ctx, tx, &models.Order{
			UserId: user.UUID,
			Amount: totalAmount,
			Date:   time.Now(),
			Status: constants.Pending,
			IsPaid: false,
		})
		if txErr != nil {
			return txErr
		}

		for _, fieldID := range request.FieldScheduleIDs {
			uuidParsed := uuid.MustParse(fieldID)
			orderFieldSchedule = append(orderFieldSchedule, models.OrderField{
				OrderID:         order.ID,
				FieldScheduleID: uuidParsed,
			})
		}

		txErr = os.repository.GetOrderField().Create(ctx, tx, orderFieldSchedule)
		if txErr != nil {
			return txErr
		}

		txErr = os.repository.GetOrderHistory().Create(ctx, tx, &dto.OrderHistoryRequest{
			Status:  constants.Pending.GetStatusString(),
			OrderID: order.ID,
		})
		if txErr != nil {
			return txErr
		}

		expiredAt := time.Now().Add(1 * time.Hour)
		description := fmt.Sprintf("Pembayaran Sewa %s", field.FieldName)
		paymentResponse, txErr = os.client.GetPayment().CreatePaymentLink(ctx, &dto.PaymentRequest{
			OrderID:     order.UUID,
			ExpiredAt:   expiredAt,
			Amount:      totalAmount,
			Description: description,
			CustomerDetail: dto.CustomerDetail{
				Name:  user.Name,
				Email: user.Email,
				Phone: user.PhoneNumber,
			},
			ItemDetails: []dto.ItemDetails{
				{
					ID:       uuid.New(),
					Name:     description,
					Amount:   totalAmount,
					Quantity: 1,
				},
			},
		})
		if txErr != nil {
			return txErr
		}

		txErr = os.repository.GetOrder().Update(ctx, tx, &models.Order{
			PaymentID: paymentResponse.UUID,
		}, order.UUID)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	response := dto.OrderResponse{
		UUID:        order.UUID,
		Code:        order.Code,
		Username:    user.Name,
		Amount:      order.Amount,
		Status:      order.Status.GetStatusString(),
		OrderDate:   order.Date,
		PaymentLink: paymentResponse.PaymentLink,
		CreatedAt:   *order.CreatedAt,
		UpdatedAt:   *order.UpdatedAt,
	}

	return &response, nil
}

func (os *OrderService) HandlePayment(context.Context, *dto.PaymentData) error {
	panic("implement me")
}

func (os *OrderService) mapPaymentStatusToOrder(request *dto.PaymentData) (constants.OrderStatus, *models.Order) {
	var (
		status constants.OrderStatus
		order  *models.Order
	)

	switch request.Status {
	case constants.SettlementPaymentStatus:
		status = constants.PaymentSuccess
		order = &models.Order{
			IsPaid:    true,
			PaymentID: request.PaymentID,
			PaidAt:    request.PaidAt,
			Status:    status,
		}
	case constants.ExpirePaymentStatus:
		status = constants.Expired
		order = &models.Order{
			IsPaid:    false,
			PaymentID: request.PaymentID,
			Status:    status,
		}
	case constants.PendingPaymentStatus:
		status = constants.PendingPayment
		order = &models.Order{
			IsPaid:    false,
			PaymentID: request.PaymentID,
			Status:    status,
		}
	}

	return status, order
}
