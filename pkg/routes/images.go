package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Images(g *echo.Group) {
	g.POST("/upload", controllers.UploadImage)
}