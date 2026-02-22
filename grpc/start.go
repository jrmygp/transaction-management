package grpc

import (
	"log"
	"net"

	"google.golang.org/grpc"

	transactionpb "github.com/jrmygp/contracts/proto/transactionpb"
	"github.com/jrmygp/transaction-management/services/order"
)

func StartGRPCServer(orderService order.OrderService) {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	transactionpb.RegisterOrderServiceServer(
		grpcServer,
		&OrderServer{Service: orderService},
	)

	log.Println("Transaction gRPC server running on :50052")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
