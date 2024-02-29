package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/auth"
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Cart(g *echo.Group) {
	g.GET("", controllers.GetCartItems, auth.WithAuth)

	g.POST("", controllers.CrerateCartItem, auth.WithAuth)

	g.PATCH("/:id", controllers.UpdateCartItem, auth.WithAuth)

	g.DELETE("/:id", controllers.DeleteCartItem, auth.WithAuth)
}
