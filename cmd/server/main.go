package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/config"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/pkg/repository"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/pkg/service"

	"github.com/joho/godotenv"

	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error is occurred  on .env file please check")
	}
	// Initialize the database
	dbConfig := config.NewConfig()
	rbConfig := config.NewRabbitMqCofig()
	// Initialize the repositories
	cashDrawerRepo := repository.NewPosCashDrawerRepository(dbConfig.SQLDB, dbConfig.RedisDB)
	customerRepo := repository.NewPosCustomerRepository(dbConfig.SQLDB, dbConfig.RedisDB)
	invoiceRepo := repository.NewPosInvoiceRepository(dbConfig.SQLDB, dbConfig.RedisDB)
	onlinePaymentRepo := repository.NewPosOnlinePaymentRepository(dbConfig.SQLDB, dbConfig.RedisDB)
	paymentMethodRepo := repository.NewPosPaymentMethodRepository(dbConfig.SQLDB, dbConfig.RedisDB)
	returnRepo := repository.NewPosReturnRepository(dbConfig.SQLDB, dbConfig.RedisDB)
	saleRepo := repository.NewPosSaleRepository(dbConfig.SQLDB, dbConfig.RedisDB)

	// Initialize the services
	cashDrawerSvc := service.NewPosCashDrawerServiceServer(cashDrawerRepo)
	customerSvc := service.NewPosCustomerService(customerRepo)
	invoiceSvc := service.NewPosInvoiceService(invoiceRepo)
	onlinePaymentSvc := service.NewPosOnlinePaymentService(onlinePaymentRepo)
	paymentMethodSvc := service.NewPosPaymentMethodService(paymentMethodRepo)
	returnSvc := service.NewPosReturnService(returnRepo)
	saleSvc := service.NewPosSaleService(saleRepo, invoiceRepo, cashDrawerRepo, onlinePaymentRepo, paymentMethodRepo, customerRepo, rbConfig.RabbitMQConn)

	// Create a gRPC server
	s := grpc.NewServer()

	// Register the services with the gRPC server
	pb.RegisterPosCashDrawerServiceServer(s, cashDrawerSvc)
	pb.RegisterPosCustomerServiceServer(s, customerSvc)
	pb.RegisterPosInvoiceServiceServer(s, invoiceSvc)
	pb.RegisterPosOnlinePaymentServiceServer(s, onlinePaymentSvc)
	pb.RegisterPosPaymentMethodServiceServer(s, paymentMethodSvc)
	pb.RegisterPosReturnServiceServer(s, returnSvc)
	pb.RegisterPosSaleServiceServer(s, saleSvc)

	// Start the gRPC server
	serverPort := os.Getenv("SERVER_PORT")
	lis, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
