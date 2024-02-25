package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Cart(g *echo.Group) {
	g.GET("/:userID", controllers.GetCartItems)

	g.POST("", controllers.CrerateCartItem)

	g.PATCH("/:id", controllers.UpdateCartItem)

	g.DELETE("/:id", controllers.DeleteCartItem)
}
