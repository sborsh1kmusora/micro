package v1

import (
	"github.com/sborsh1kmusora/micro/payment/internal/service"
	paymentV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewAPI(paymentService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
