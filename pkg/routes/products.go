package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Products(g *echo.Group) {
	g.GET("", controllers.GetProducts)
	g.GET("/:id", controllers.GetProduct)

	g.POST("", controllers.CreateProduct)
	g.POST("/:id/image", controllers.UploadProductImage)

	g.PUT("/:id", controllers.UpdateProduct)
	// g.PATCH("/:id", controllers.PatchProduct)

	g.DELETE("/:id", controllers.DeleteProduct)
	g.DELETE("/:id/image", controllers.DeleteProductImage)
}
