package constrollers

import (
	"net/http"
	errConst "order-service/common/error"
	"order-service/common/response"
	"order-service/domain/dto"
	"order-service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OrderController struct {
	service services.IRegistryService
}

type IOrderController interface {
	GetAllWithPagination(*gin.Context)
	GetByUUID(*gin.Context)
	GetOrderByUserID(*gin.Context)
	Create(*gin.Context)
}

func NewOrderController(service services.IRegistryService) IOrderController {
	return &OrderController{service: service}
}

func (oc *OrderController) GetAllWithPagination(ctx *gin.Context) {
	var params dto.OrderRequestParam
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	if err = validate.Struct(params); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errConst.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResponse{
			Err:     err,
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     ctx,
		})
		return
	}

	result, err := oc.service.GetOrder().GetAllWithPagination(ctx, &params)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: result,
		Gin:  ctx,
	})
}

func (oc *OrderController) GetByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	result, err := oc.service.GetOrder().GetByUUID(ctx, uuid)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: result,
		Gin:  ctx,
	})
}

func (oc *OrderController) GetOrderByUserID(ctx *gin.Context) {
	result, err := oc.service.GetOrder().GetOrderByUserID(ctx.Request.Context())
	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: result,
		Gin:  ctx,
	})
}

func (oc *OrderController) Create(*gin.Context) {}
