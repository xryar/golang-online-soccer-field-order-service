package services

import (
	"order-service/clients"
	"order-service/repositories"
	services "order-service/services/order"
)

type Registry struct {
	repository repositories.IRegistryRepository
	client     clients.IRegistryClient
}

type IRegistryService interface {
	GetOrder() services.IOrderService
}

func NewRegistryService(repository repositories.IRegistryRepository, client clients.IRegistryClient) IRegistryService {
	return &Registry{
		repository: repository,
		client:     client,
	}
}

func (r *Registry) GetOrder() services.IOrderService {
	return services.NewOrderService(r.repository, r.client)
}
