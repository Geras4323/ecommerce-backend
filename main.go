package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/geras4323/ecommerce-backend/pkg/cloud"
	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	app := echo.New()

	if errDotenv := godotenv.Load(); errDotenv != nil {
		app.Logger.Fatal("Error loading environment variables")
	}

	database.ConnectGorm()
	cloud.ConnectCloudinary()

	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "https://misideaspintadas.com.ar"},
		AllowHeaders:     []string{},
		AllowCredentials: true,
	}))

	app.Pre(middleware.RemoveTrailingSlash())
	routes.SetupRoutes(app)

	_, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.Logger.Fatal(app.Start(":1323"))
}
