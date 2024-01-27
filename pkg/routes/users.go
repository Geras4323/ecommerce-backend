package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Users(g *echo.Group) {
	g.GET("", controllers.GetUsers)
	g.GET("/:id", controllers.GetUser)

	g.POST("", controllers.CreateUser)

	g.PUT("/:id", controllers.UpdateUser)
	// g.PATCH("/:id", controllers.PatchUser)

	g.DELETE("/:id", controllers.DeleteUser)
}
