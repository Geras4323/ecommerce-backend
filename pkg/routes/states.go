package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func States(g *echo.Group) {
	g.POST("/vacation/set", controllers.SetVacation)

	g.GET("/vacation", controllers.GetVacation)
	g.PATCH("/vacation", controllers.UpdateVacation)
}
