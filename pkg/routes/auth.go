package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/auth"
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Auth(g *echo.Group) {
	g.GET("/session", controllers.GetSession, auth.WithAuth)

	g.POST("/login", controllers.Login)
	g.POST("/logout", controllers.Logout, auth.WithAuth)

	g.POST("/signup", controllers.Signup)
	g.POST("/signup/verify/:token", controllers.VerifyEmail)
	g.POST("/signup/verify/restart", controllers.RestarEmailVerification, auth.WithAuth)

	g.POST("/recovery", controllers.RecoverPassword)
	g.POST("/change-password", controllers.ChangePassword)
}
