package routes

import "github.com/labstack/echo/v4"

func SetupRoutes(app *echo.Echo) {
	api := app.Group("/api/v1")

	Categories(api.Group("/categories"))
	Suppliers(api.Group("/suppliers"))
	Users(api.Group("/users"))
	Cart(api.Group("/cart"))
	Auth(api.Group("/auth"))
	Products(api.Group("/products"))
	Orders(api.Group("/orders"))
	Payments(api.Group("/payments"))
	Images(api.Group("/images"))

	Email(api.Group("/email"))
}
