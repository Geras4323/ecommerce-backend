package routes

import (
	"github.com/geras4323/ecommerce-backend/pkg/auth"
	"github.com/geras4323/ecommerce-backend/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func Images(g *echo.Group) {
	g.GET("", controllers.GetImages)

	g.POST("/test", controllers.ImagesTest)
	g.POST("/upload", controllers.UploadImage, auth.WithAuth, auth.CheckAdmin)

	g.POST("/upload/pdf", controllers.UploadPDF, auth.WithAuth, auth.CheckAdmin)
}
