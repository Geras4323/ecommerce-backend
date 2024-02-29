package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/auth"
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Orders(g *echo.Group) {
	// controller, require user to be logged, check role
	// g.GET("", controllers.GetOrders, echojwt.WithConfig(config), auth.CheckRole("admin", "customer"))

	g.GET("", controllers.GetOrders, auth.WithAuth, auth.CheckAdmin)
	g.GET("/my-orders", controllers.GetOrdersByUser, auth.WithAuth)
	g.GET("/:id", controllers.GetOrder, auth.WithAuth)

	g.POST("", controllers.CreateOrder, auth.WithAuth)

	g.PATCH("/:id", controllers.UpdateOrder, auth.WithAuth, auth.CheckAdmin)

	g.DELETE("/:id", controllers.DeleteOrder, auth.WithAuth, auth.CheckAdmin)
}
