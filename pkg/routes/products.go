package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/auth"
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Products(g *echo.Group) {
	g.GET("", controllers.GetProducts)
	g.GET("/:id", controllers.GetProduct, auth.WithAuth, auth.CheckAdmin)

	g.POST("", controllers.CreateProduct, auth.WithAuth, auth.CheckAdmin)
	g.POST("/:id/image", controllers.UploadProductImage, auth.WithAuth, auth.CheckAdmin)

	g.PUT("/:id", controllers.UpdateProduct, auth.WithAuth, auth.CheckAdmin)

	g.DELETE("/:id", controllers.DeleteProduct, auth.WithAuth, auth.CheckAdmin)
	g.DELETE("/:id/image", controllers.DeleteProductImage, auth.WithAuth, auth.CheckAdmin)
}
