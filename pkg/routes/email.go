package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Email(g *echo.Group) {
	g.GET("", controllers.GetEmail)
}