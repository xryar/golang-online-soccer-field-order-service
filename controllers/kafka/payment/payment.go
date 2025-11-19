package kafka

import (
	"context"
	"encoding/json"
	"order-service/common/util"
	"order-service/domain/dto"
	"order-service/services"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

const PaymentTopic = "payment-service-callback"

type PaymentKafka struct {
	service services.IRegistryService
}

type IPaymentKafka interface {
	HandlePayment(context.Context, *sarama.ConsumerMessage) error
}

func NewPaymentKafka(service services.IRegistryService) IPaymentKafka {
	return &PaymentKafka{service: service}
}

func (pk *PaymentKafka) HandlePayment(ctx context.Context, message *sarama.ConsumerMessage) error {
	defer util.Recover()
	var body dto.PaymentContent
	err := json.Unmarshal(message.Value, &body)
	if err != nil {
		logrus.Errorf("failed to unmarshal message: %v", err)
		return err
	}

	data := body.Body.Data
	err = pk.service.GetOrder().HandlePayment(ctx, &data)
	if err != nil {
		logrus.Errorf("failed to handle payment: %v", err)
		return err
	}

	logrus.Infof("success handle payment")
	return nil
}
