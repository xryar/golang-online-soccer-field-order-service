package controllers

import (
	controllers "order-service/controllers/http/order"
	"order-service/services"
)

type Registry struct {
	service services.IRegistryService
}

type IRegistryController interface {
	GetOrder() controllers.IOrderController
}

func NewRegistryController(service services.IRegistryService) IRegistryController {
	return &Registry{service: service}
}

func (r *Registry) GetOrder() controllers.IOrderController {
	return controllers.NewOrderController(r.service)
}
