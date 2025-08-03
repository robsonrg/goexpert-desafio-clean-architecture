package service

import (
	"context"

	"github.com/robsonrg/goexpert-desafio-clean-architecture/internal/infra/grpc/pb"
	"github.com/robsonrg/goexpert-desafio-clean-architecture/internal/usecase/orders"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase orders.CreateOrderUseCase
	ListOrderUseCase   orders.ListOrderUseCase
}

func NewOrderService(
	createOrderUseCase orders.CreateOrderUseCase,
	listOrderUseCase orders.ListOrderUseCase,
) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrderUseCase:   listOrderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := orders.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	listOrdersResp, err := s.ListOrderUseCase.GetOrders()
	if err != nil {
		return nil, err
	}

	resp := &pb.ListOrdersResponse{}
	for _, order := range listOrdersResp.Orders {
		resp.Orders = append(resp.Orders, &pb.Order{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		})
	}
	return resp, nil
}
