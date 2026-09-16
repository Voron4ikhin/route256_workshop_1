package grpc

import (
	"context"
	portin "route256/loms/internal/port/in"
	lomsv1 "route256/loms/pkg/api/loms/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Dependencies struct {
	OrderCreator   portin.OrderCreator
	OrderPayer     portin.OrderPayer
	OrderCanceler  portin.OrderCanceler
	OrderInformant portin.OrderInformant
	StockInformant portin.StockInformant
}

type Server struct {
	lomsv1.UnimplementedLOMSServiceServer
	deps Dependencies
}

func NewServer(deps Dependencies) *Server {
	return &Server{deps: deps}
}

func Register(s *grpc.Server, deps Dependencies) {
	lomsv1.RegisterLOMSServiceServer(s, NewServer(deps))
}

func (s *Server) CreateOrder(ctx context.Context, req *lomsv1.CreateOrderRequest) (*lomsv1.CreateOrderResponse, error) {
	if req.GetUser() <= 0 {
		return nil, status.Error(codes.InvalidArgument, ErrIncorrectUser.Error())
	}
	if len(req.GetItems()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ErrIncorrectQuantity.Error())
	}
	for _, item := range req.GetItems() {
		if item.GetSku() == 0 {
			return nil, status.Error(codes.InvalidArgument, ErrIncorrectSKU.Error())
		}
		if item.GetCount() <= 0 {
			return nil, status.Error(codes.InvalidArgument, ErrIncorrectProductCount.Error())
		}
	}

	orderID, err := s.deps.OrderCreator.CreateOrder(ctx, req.GetUser(), itemsToDomain(req.GetItems()))
	if err != nil {
		return nil, status.Error(classifyBusinessError(err), err.Error())
	}

	return &lomsv1.CreateOrderResponse{OrderId: orderID}, nil
}

func (s *Server) GetOrderInfo(ctx context.Context, req *lomsv1.GetOrderInfoRequest) (*lomsv1.GetOrderInfoResponse, error) {
	if req.GetOrderId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, ErrIncorrectOrderID.Error())
	}

	order, err := s.deps.OrderInformant.GetOrderInfo(ctx, req.GetOrderId())
	if err != nil {
		return nil, status.Error(classifyBusinessError(err), err.Error())
	}

	return &lomsv1.GetOrderInfoResponse{
		Status: string(order.Status),
		User:   order.UserID,
		Items:  itemsFromDomain(order.Items),
	}, nil
}

func (s *Server) PayOrder(ctx context.Context, req *lomsv1.PayOrderRequest) (*lomsv1.PayOrderResponse, error) {
	if req.GetOrderId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, ErrIncorrectOrderID.Error())
	}

	if err := s.deps.OrderPayer.PayOrder(ctx, req.GetOrderId()); err != nil {
		return nil, status.Error(classifyBusinessError(err), err.Error())
	}

	return &lomsv1.PayOrderResponse{}, nil
}

func (s *Server) CancelOrder(ctx context.Context, req *lomsv1.CancelOrderRequest) (*lomsv1.CancelOrderResponse, error) {
	if req.GetOrderId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, ErrIncorrectOrderID.Error())
	}

	if err := s.deps.OrderCanceler.CancelOrder(ctx, req.GetOrderId()); err != nil {
		return nil, status.Error(classifyBusinessError(err), err.Error())
	}

	return &lomsv1.CancelOrderResponse{}, nil
}

func (s *Server) GetStocks(ctx context.Context, req *lomsv1.GetStocksRequest) (*lomsv1.GetStocksResponse, error) {
	if req.GetSku() == 0 {
		return nil, status.Error(codes.InvalidArgument, ErrIncorrectSKU.Error())
	}
	count := s.deps.StockInformant.GetStocks(ctx, req.GetSku())

	return &lomsv1.GetStocksResponse{Count: count}, nil
}
