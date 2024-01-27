package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Suppliers(g *echo.Group) {
	g.GET("", controllers.GetSuppliers)
	g.GET("/:id", controllers.GetSupplier)

	g.POST("", controllers.CreateSupplier)

	g.PUT("/:id", controllers.UpdateSupplier)
	// g.PATCH("/:id", controllers.PatchSupplier)

	g.DELETE("/:id", controllers.DeleteSupplier)
}
