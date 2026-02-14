package v1

import (
	def "github.com/sborsh1kmusora/micro/order/internal/client/grpc"
	paymentV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	genClient paymentV1.PaymentServiceClient
}

func NewClient(genClient paymentV1.PaymentServiceClient) *client {
	return &client{
		genClient: genClient,
	}
}
