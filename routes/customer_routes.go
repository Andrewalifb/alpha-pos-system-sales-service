package routes

import (
	"github.com/Andrewalifb/alpha-pos-system-sales-service/api/controller"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/api/midlleware"
	"github.com/gin-gonic/gin"
)

func PosCustomerRoutes(r *gin.Engine, posCustomerController controller.PosCustomerController) {
	routes := r.Group("/api")

	// Apply the JWT middleware to all routes in this group
	routes.Use(midlleware.JWTAuthMiddleware())

	routesV1 := routes.Group("/v1/customers")
	// Create New PosCustomer
	routesV1.POST("/pos_customer", posCustomerController.HandleCreatePosCustomerRequest)
	// Get PosCustomer by ID
	routesV1.GET("/pos_customer/:id", posCustomerController.HandleReadPosCustomerRequest)
	// Update Existing PosCustomer
	routesV1.PUT("/pos_customer/:id", posCustomerController.HandleUpdatePosCustomerRequest)
	// Delete PosCustomer
	routesV1.DELETE("/pos_customer/:id", posCustomerController.HandleDeletePosCustomerRequest)
	// Get All PosCustomers
	routesV1.GET("/pos_customers", posCustomerController.HandleReadAllPosCustomersRequest)
}
