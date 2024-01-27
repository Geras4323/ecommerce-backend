package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Products(g *echo.Group) {
	g.GET("", controllers.GetProducts)
	g.GET("/:id", controllers.GetProduct)

	g.POST("", controllers.CreateProduct)

	g.PUT("/:id", controllers.UpdateProduct)
	// g.PATCH("/:id", controllers.PatchProduct)

	g.DELETE("/:id", controllers.DeleteProduct)
}
