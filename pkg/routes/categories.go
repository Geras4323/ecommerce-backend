package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/auth"
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Categories(g *echo.Group) {
	g.GET("", controllers.GetCategories)
	g.GET("/:id", controllers.GetCategory, auth.WithAuth)

	g.POST("", controllers.CreateCategory, auth.WithAuth, auth.CheckAdmin)
	g.POST("/:id/image", controllers.UploadCategoryImage, auth.WithAuth, auth.CheckAdmin)

	g.PUT("/:id", controllers.UpdateCategory, auth.WithAuth, auth.CheckAdmin)

	g.DELETE("/:id", controllers.DeleteCategory, auth.WithAuth, auth.CheckAdmin)
	g.DELETE("/:id/image", controllers.DeleteCategoryImage, auth.WithAuth, auth.CheckAdmin)
}
