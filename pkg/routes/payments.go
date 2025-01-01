package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Payments(g *echo.Group) {
	g.GET("", controllers.GetPayments)
	g.GET("/:id", controllers.GetPayment)

	g.POST("/:orderID", controllers.CreatePayment)
	g.POST("/mercadopago/add", controllers.AddMPPayment)
	g.POST("/mercadopago/:id/end", controllers.EndMPPayment)

	g.PUT("/:id", controllers.UpdatePayment)

	g.DELETE("/:id", controllers.DeletePayment)
}
