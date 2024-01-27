package controllers

import (
	"fmt"
	"net/http"

	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GET /api/v1/orders //////////////////////////////////////////////////////
func GetOrders(c echo.Context) error {
	orders := make([]models.Order, 0)

	if err := database.Gorm.Preload("Payments").Find(&orders).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, orders)
}

// GET /api/v1/orders/my-orders
func GetOrdersByUser(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	UserID := claims["Sub"]

	var orders []models.Order

	if err := database.Gorm.Where("user_id = ?", UserID).Find(&orders).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, orders)
}

// GET /api/v1/orders/:id
func GetOrder(c echo.Context) error {
	orderID := c.Param("id")

	var order models.Order
	if err := database.Gorm.Preload("Payments").First(&order, orderID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, order)
}

// POST /api/v1/orders //////////////////////////////////////////////////////
func CreateOrder(c echo.Context) error {
	body := models.CreateOrder{}
	c.Bind(&body)

	order := models.Order{
		UserID: body.UserID,
		Total:  body.Total,
	}

	if err := database.Gorm.Create(&order).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, order)
}

// POST /api/v1/orders/add-product
func AddProduct(c echo.Context) error {
	var body models.AddProduct
	c.Bind(&body)

	if err := database.Gorm.Create(&body).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, body)
}

// PUT /api/v1/orders/:id //////////////////////////////////////////////////////
func UpdateOrder(c echo.Context) error {
	orderID := c.Param("id")

	var newOrder models.UpdateOrder
	c.Bind(&newOrder)

	var oldOrder models.Order
	if err := database.Gorm.First(&oldOrder, orderID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	oldOrder.UserID = newOrder.UserID
	oldOrder.Total = newOrder.Total

	if err := database.Gorm.Where("id = ?", orderID).Save(&oldOrder).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, oldOrder)
}

// PATCH /api/v1/orders/:id
// func PatchOrder(c echo.Context) error {

// 	return nil
// }

// DELETE /api/v1/orders/:id //////////////////////////////////////////////////////
func DeleteOrder(c echo.Context) error {
	orderID := c.Param("id")

	if err := database.Gorm.First(&models.Order{}, orderID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	if err := database.Gorm.Delete(&models.Order{}, orderID).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{"id": orderID, "message": "Order deleted successfully"})
}

// DELETE /api/v1/orders/:id/remove-product/:productID
func RemoveProduct(c echo.Context) error {
	orderID := c.Param("id")
	productID := c.Param("productID")

	if err := database.Gorm.Where("order_id = ? AND product_id = ?", orderID, productID).Delete(&models.OrderProduct{}).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{"orderID": orderID, "productID": productID, "message": fmt.Sprintf("Product %s removed successfully from order %s", productID, orderID)})
}
