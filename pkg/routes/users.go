package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/auth"
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Users(g *echo.Group) {
	g.GET("", controllers.GetUsers, auth.WithAuth, auth.CheckAdmin)
	g.GET("/:id", controllers.GetUser, auth.WithAuth, auth.CheckAdmin)

	g.POST("", controllers.CreateUser, auth.WithAuth, auth.CheckAdmin)

	g.PUT("/:id", controllers.UpdateUser, auth.WithAuth, auth.CheckAdmin)
	// g.PATCH("/:id", controllers.PatchUser)

	g.DELETE("/:id", controllers.DeleteUser, auth.WithAuth, auth.CheckAdmin)
}
