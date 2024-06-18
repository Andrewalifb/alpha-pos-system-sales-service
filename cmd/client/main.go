package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Andrewalifb/alpha-pos-system-sales-service/api/controller"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/routes"

	pb "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error is occurred  on .env file please check")
	}

	clientPort := os.Getenv("CLIENT_PORT")
	serverPort := os.Getenv("SERVER_PORT")

	addr := fmt.Sprintf("localhost:%s", serverPort)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create gRPC clients for each service
	cashDrawerClient := pb.NewPosCashDrawerServiceClient(conn)
	customerClient := pb.NewPosCustomerServiceClient(conn)
	invoiceClient := pb.NewPosInvoiceServiceClient(conn)
	onlinePaymentClient := pb.NewPosOnlinePaymentServiceClient(conn)
	paymentMethodClient := pb.NewPosPaymentMethodServiceClient(conn)
	returnClient := pb.NewPosReturnServiceClient(conn)
	saleClient := pb.NewPosSaleServiceClient(conn)

	// Initialize the controllers with the gRPC clients
	cashDrawerCtrl := controller.NewPosCashDrawerController(cashDrawerClient)
	customerCtrl := controller.NewPosCustomerController(customerClient)
	invoiceCtrl := controller.NewPosInvoiceController(invoiceClient)
	onlinePaymentCtrl := controller.NewPosOnlinePaymentController(onlinePaymentClient)
	paymentMethodCtrl := controller.NewPosPaymentMethodController(paymentMethodClient)
	returnCtrl := controller.NewPosReturnController(returnClient)
	saleCtrl := controller.NewPosSaleController(saleClient)

	// Create a new router
	r := gin.Default()

	// Define your routes
	routes.PosCashDrawerRoutes(r, cashDrawerCtrl)
	routes.PosCustomerRoutes(r, customerCtrl)
	routes.PosInvoiceRoutes(r, invoiceCtrl)
	routes.PosOnlinePaymentRoutes(r, onlinePaymentCtrl)
	routes.PosPaymentMethodRoutes(r, paymentMethodCtrl)
	routes.PosReturnRoutes(r, returnCtrl)
	routes.PosSaleRoutes(r, saleCtrl)

	// Start the server
	r.Run(":" + clientPort)
}
