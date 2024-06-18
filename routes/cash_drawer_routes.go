package routes

import (
	"github.com/Andrewalifb/alpha-pos-system-sales-service/api/controller"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/api/midlleware"

	"github.com/gin-gonic/gin"
)

func PosCashDrawerRoutes(r *gin.Engine, posCashDrawerController controller.PosCashDrawerController) {
	routes := r.Group("/api")

	// Apply the JWT middleware to all routes in this group
	routes.Use(midlleware.JWTAuthMiddleware())

	routesV1 := routes.Group("/v1/cash-drawers")
	// Create New PosCashDrawer
	routesV1.POST("/pos_cash_drawer", posCashDrawerController.HandleCreatePosCashDrawerRequest)
	// Get PosCashDrawer by ID
	routesV1.GET("/pos_cash_drawer/:id", posCashDrawerController.HandleReadPosCashDrawerRequest)
	// Update Existing PosCashDrawer
	routesV1.PUT("/pos_cash_drawer/:id", posCashDrawerController.HandleUpdatePosCashDrawerRequest)
	// Delete PosCashDrawer
	routesV1.DELETE("/pos_cash_drawer/:id", posCashDrawerController.HandleDeletePosCashDrawerRequest)
	// Get All PosCashDrawers
	routesV1.GET("/pos_cash_drawers", posCashDrawerController.HandleReadAllPosCashDrawersRequest)
}
