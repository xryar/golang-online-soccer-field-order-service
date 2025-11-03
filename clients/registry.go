package clients

import (
	"order-service/clients/config"
	fieldConfig "order-service/clients/field"
	paymentConfig "order-service/clients/payment"
	userConfig "order-service/clients/user"
	appConfig "order-service/config"
)

type RegistryClient struct{}

type IRegistryClient interface {
	GetUser() userConfig.IUserClient
	GetPayment() paymentConfig.IPaymentClient
	GetField() fieldConfig.IFieldClient
}

func NewRegistryClient() IRegistryClient {
	return &RegistryClient{}
}

func (rc *RegistryClient) GetUser() userConfig.IUserClient {
	return userConfig.NewUserClient(
		config.NewClientConfig(
			config.WithBaseURL(appConfig.Config.InternalService.User.Host),
			config.WithSignatureKey(appConfig.Config.InternalService.User.SignatureKey),
		),
	)
}

func (rc *RegistryClient) GetPayment() paymentConfig.IPaymentClient {
	return paymentConfig.NewPaymentClient(
		config.NewClientConfig(
			config.WithBaseURL(appConfig.Config.InternalService.User.Host),
			config.WithSignatureKey(appConfig.Config.InternalService.User.SignatureKey),
		),
	)
}

func (rc *RegistryClient) GetField() fieldConfig.IFieldClient {
	return fieldConfig.NewFieldClient(
		config.NewClientConfig(
			config.WithBaseURL(appConfig.Config.InternalService.User.Host),
			config.WithSignatureKey(appConfig.Config.InternalService.User.SignatureKey),
		),
	)
}
