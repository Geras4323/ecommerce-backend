package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/auth"
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Suppliers(g *echo.Group) {
	g.GET("", controllers.GetSuppliers, auth.WithAuth, auth.CheckAdmin)
	g.GET("/:id", controllers.GetSupplier, auth.WithAuth, auth.CheckAdmin)

	g.POST("", controllers.CreateSupplier, auth.WithAuth, auth.CheckAdmin)

	g.PUT("/:id", controllers.UpdateSupplier, auth.WithAuth, auth.CheckAdmin)

	g.DELETE("/:id", controllers.DeleteSupplier, auth.WithAuth, auth.CheckAdmin)
}
