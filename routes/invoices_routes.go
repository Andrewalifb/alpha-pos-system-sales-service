package routes

import (
	"github.com/Andrewalifb/alpha-pos-system-sales-service/api/controller"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/api/midlleware"
	"github.com/gin-gonic/gin"
)

func PosInvoiceRoutes(r *gin.Engine, posInvoiceController controller.PosInvoiceController) {
	routes := r.Group("/api")

	// Apply the JWT middleware to all routes in this group
	routes.Use(midlleware.JWTAuthMiddleware())

	routesV1 := routes.Group("/v1/invoices")
	// Create New PosInvoice
	routesV1.POST("/pos_invoice", posInvoiceController.HandleCreatePosInvoiceRequest)
	// Get PosInvoice by ID
	routesV1.GET("/pos_invoice/:id", posInvoiceController.HandleReadPosInvoiceRequest)
	// Update Existing PosInvoice
	routesV1.PUT("/pos_invoice/:id", posInvoiceController.HandleUpdatePosInvoiceRequest)
	// Delete PosInvoice
	routesV1.DELETE("/pos_invoice/:id", posInvoiceController.HandleDeletePosInvoiceRequest)
	// Get All PosInvoices
	routesV1.GET("/pos_invoices", posInvoiceController.HandleReadAllPosInvoicesRequest)
}
