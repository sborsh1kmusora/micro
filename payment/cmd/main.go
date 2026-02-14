package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentV1API "github.com/sborsh1kmusora/micro/payment/internal/api/payment/v1"
	"github.com/sborsh1kmusora/micro/payment/internal/interceptor"
	paymentService "github.com/sborsh1kmusora/micro/payment/internal/service/payment"
	paymentV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggerInterceptor()),
	)

	svc := paymentService.NewService()
	api := paymentV1API.NewAPI(svc)

	paymentV1.RegisterPaymentServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("grpc server listening at %d", grpcPort)
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Printf("Shutting down gRPC server...")
	s.GracefulStop()
	log.Printf("Server stopped gracefully")
}
