package graph

import "github.com/robsonrg/goexpert-desafio-clean-architecture/internal/usecase/orders"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CreateOrderUseCase orders.CreateOrderUseCase
	ListOrderUseCase   orders.ListOrderUseCase
}
