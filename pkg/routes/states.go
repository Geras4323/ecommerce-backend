package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func States(g *echo.Group) {
	g.GET("/vacation", controllers.GetVacation)
	g.PATCH("/vacation", controllers.UpdateVacation)

	g.GET("/mercadopago", controllers.GetMPPayments)
	g.PATCH("/mercadopago", controllers.UpdateMPPayments)
}
