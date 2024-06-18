package routes

import (
	"github.com/Andrewalifb/alpha-pos-system-sales-service/api/controller"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/api/midlleware"
	"github.com/gin-gonic/gin"
)

func PosSaleRoutes(r *gin.Engine, posSaleController controller.PosSaleController) {
	routes := r.Group("/api")

	// Apply the JWT middleware to all routes in this group
	routes.Use(midlleware.JWTAuthMiddleware())

	routesV1 := routes.Group("/v1/sales")
	// Create New PosSale
	routesV1.POST("/pos_sale", posSaleController.HandleCreatePosSaleRequest)
	// Get PosSale by ID
	routesV1.GET("/pos_sale/:id", posSaleController.HandleReadPosSaleRequest)
	// Update Existing PosSale
	routesV1.PUT("/pos_sale/:id", posSaleController.HandleUpdatePosSaleRequest)
	// Delete PosSale
	routesV1.DELETE("/pos_sale/:id", posSaleController.HandleDeletePosSaleRequest)
	// Get All PosSales
	routesV1.GET("/pos_sales", posSaleController.HandleReadAllPosSalesRequest)
}
