package clients

import (
	orderConfig "order-service/clients/config"
	userConfig "order-service/clients/user"
	"order-service/config"
)

type RegistryClient struct{}

type IRegistryClient interface {
	GetUser() userConfig.IUserClient
}

func NewRegistryClient() IRegistryClient {
	return &RegistryClient{}
}

func (rc *RegistryClient) GetUser() userConfig.IUserClient {
	return userConfig.NewUserClient(
		orderConfig.NewClientConfig(
			orderConfig.WithBaseURL(config.Config.InternalService.User.Host),
			orderConfig.WithSignatureKey(config.Config.InternalService.User.SignatureKey),
		),
	)
}
