package routes

import (
	"github.com/Andrewalifb/alpha-pos-system-sales-service/api/controller"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/api/midlleware"
	"github.com/gin-gonic/gin"
)

func PosReturnRoutes(r *gin.Engine, posReturnController controller.PosReturnController) {
	routes := r.Group("/api")

	// Apply the JWT middleware to all routes in this group
	routes.Use(midlleware.JWTAuthMiddleware())

	routesV1 := routes.Group("/v1/returns")
	// Create New PosReturn
	routesV1.POST("/pos_return", posReturnController.HandleCreatePosReturnRequest)
	// Get PosReturn by ID
	routesV1.GET("/pos_return/:id", posReturnController.HandleReadPosReturnRequest)
	// Update Existing PosReturn
	routesV1.PUT("/pos_return/:id", posReturnController.HandleUpdatePosReturnRequest)
	// Delete PosReturn
	routesV1.DELETE("/pos_return/:id", posReturnController.HandleDeletePosReturnRequest)
	// Get All PosReturns
	routesV1.GET("/pos_returns", posReturnController.HandleReadAllPosReturnsRequest)
}
