package grpc

import (
	"context"

	transactionpb "github.com/jrmygp/contracts/proto/transactionpb"
	"github.com/jrmygp/transaction-management/services/order"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderServer struct {
	transactionpb.UnimplementedOrderServiceServer
	Service order.OrderService
}

func (s *OrderServer) FindByMidtransOrderID(
	ctx context.Context,
	req *transactionpb.GetOrderByMidtransRequest,
) (*transactionpb.GetOrderResponse, error) {

	orderData, err := s.Service.FindByMidtransOrderID(req.Id)
	if err != nil {
		return nil, err
	}

	var paidAt *timestamppb.Timestamp
	if orderData.PaidAt != nil {
		paidAt = timestamppb.New(*orderData.PaidAt)
	}

	return &transactionpb.GetOrderResponse{
		Id:              int32(orderData.ID),
		UserId:          int32(orderData.UserID),
		HotelId:         int32(orderData.HotelID),
		Nights:          int32(orderData.Nights),
		Status:          orderData.Status,
		Bill:            int32(orderData.Bill),
		MidtransOrderId: orderData.MidtransOrderID,
		PaymentUrl:      orderData.PaymentURL,
		PaymentMethod:   orderData.PaymentMethod,
		PaidAt:          paidAt,
	}, nil
}
