package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/geras4323/ecommerce-backend/pkg/utils"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Orders(g *echo.Group) {
	config := echojwt.Config{
		SigningKey: []byte(utils.GetEnvVar("JWT_LOGIN_SECRET")),
	}

	// controller, require user to be logged, check role
	// g.GET("", controllers.GetOrders, echojwt.WithConfig(config), auth.CheckRole("admin", "customer"))
	g.GET("", controllers.GetOrders)
	g.GET("/my-orders", controllers.GetOrdersByUser, echojwt.WithConfig(config))
	g.GET("/:id", controllers.GetOrder)

	g.POST("/user/:userID", controllers.CreateOrder)
	g.POST("/:id/add-product", controllers.AddProduct)

	g.PATCH("/:id", controllers.UpdateOrder)

	g.DELETE("/:id", controllers.DeleteOrder)
	g.DELETE("/:id/remove-product/:productID", controllers.RemoveProduct)
}
