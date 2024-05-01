package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/auth"
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Orders(g *echo.Group) {
	g.GET("", controllers.GetOrders, auth.WithAuth, auth.CheckAdmin)
	g.GET("/my-orders", controllers.GetOrdersByUser, auth.WithAuth)
	g.GET("/:id", controllers.GetOrder, auth.WithAuth)

	g.POST("", controllers.CreateOrder, auth.WithAuth)

	g.PATCH("/:id/state", controllers.UpdateOrderState, auth.WithAuth, auth.CheckAdmin)

	g.DELETE("/:id", controllers.DeleteOrder, auth.WithAuth, auth.CheckAdmin)
}
