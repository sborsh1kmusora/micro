package v1

import "github.com/sborsh1kmusora/micro/order/internal/service"

type api struct {
	orderSvc service.OrderService
}

func NewAPI(svc service.OrderService) *api {
	return &api{
		orderSvc: svc,
	}
}
