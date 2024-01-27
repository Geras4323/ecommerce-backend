package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Payments(g *echo.Group) {
	g.GET("", controllers.GetPayments)
	g.GET("/:id", controllers.GetPayment)

	g.POST("", controllers.CreatePayment)

	g.PUT("/:id", controllers.UpdatePayment)
	// g.PATCH("/:id", controllers.PatchPayment)

	g.DELETE("/:id", controllers.DeletePayment)
}
