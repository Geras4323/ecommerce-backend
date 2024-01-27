package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Categories(g *echo.Group) {
	g.GET("", controllers.GetCategories)
	g.GET("/:id", controllers.GetCategory)

	g.POST("", controllers.CreateCategory)
	g.POST("/:id/image", controllers.UploadCategoryImage)

	g.PUT("/:id", controllers.UpdateCategory)

	g.DELETE("/:id", controllers.DeleteCategory)
	g.DELETE("/:id/image", controllers.DeleteCategoryImage)
}
