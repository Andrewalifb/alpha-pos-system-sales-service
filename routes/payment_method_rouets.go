package routes

import (
	"github.com/Andrewalifb/alpha-pos-system-sales-service/api/controller"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/api/midlleware"
	"github.com/gin-gonic/gin"
)

func PosPaymentMethodRoutes(r *gin.Engine, posPaymentMethodController controller.PosPaymentMethodController) {
	routes := r.Group("/api")

	// Apply the JWT middleware to all routes in this group
	routes.Use(midlleware.JWTAuthMiddleware())

	routesV1 := routes.Group("/v1/payment-methods")
	// Create New PosPaymentMethod
	routesV1.POST("/pos_payment_method", posPaymentMethodController.HandleCreatePosPaymentMethodRequest)
	// Get PosPaymentMethod by ID
	routesV1.GET("/pos_payment_method/:id", posPaymentMethodController.HandleReadPosPaymentMethodRequest)
	// Update Existing PosPaymentMethod
	routesV1.PUT("/pos_payment_method/:id", posPaymentMethodController.HandleUpdatePosPaymentMethodRequest)
	// Delete PosPaymentMethod
	routesV1.DELETE("/pos_payment_method/:id", posPaymentMethodController.HandleDeletePosPaymentMethodRequest)
	// Get All PosPaymentMethods
	routesV1.GET("/pos_payment_methods", posPaymentMethodController.HandleReadAllPosPaymentMethodsRequest)
}
