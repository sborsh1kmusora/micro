package payment

import def "github.com/sborsh1kmusora/micro/payment/internal/service"

var _ def.PaymentService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
