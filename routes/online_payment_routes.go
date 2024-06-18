package routes

import (
	"github.com/Andrewalifb/alpha-pos-system-sales-service/api/controller"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/api/midlleware"
	"github.com/gin-gonic/gin"
)

func PosOnlinePaymentRoutes(r *gin.Engine, posOnlinePaymentController controller.PosOnlinePaymentController) {
	routes := r.Group("/api")

	// Apply the JWT middleware to all routes in this group
	routes.Use(midlleware.JWTAuthMiddleware())

	routesV1 := routes.Group("/v1/online-payments")
	// Create New PosOnlinePayment
	routesV1.POST("/pos_online_payment", posOnlinePaymentController.HandleCreatePosOnlinePaymentRequest)
	// Get PosOnlinePayment by ID
	routesV1.GET("/pos_online_payment/:id", posOnlinePaymentController.HandleReadPosOnlinePaymentRequest)
	// Update Existing PosOnlinePayment
	routesV1.PUT("/pos_online_payment/:id", posOnlinePaymentController.HandleUpdatePosOnlinePaymentRequest)
	// Delete PosOnlinePayment
	routesV1.DELETE("/pos_online_payment/:id", posOnlinePaymentController.HandleDeletePosOnlinePaymentRequest)
	// Get All PosOnlinePayments
	routesV1.GET("/pos_online_payments", posOnlinePaymentController.HandleReadAllPosOnlinePaymentsRequest)
}
