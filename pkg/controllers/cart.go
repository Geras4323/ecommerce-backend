package controllers

import (
	"net/http"

	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/labstack/echo/v4"
)

func GetCartItems(c echo.Context) error {
	userID := c.Param("userID")

	if err := database.Gorm.First(&models.User{}, userID).Error; err != nil {
		return c.String(http.StatusNotFound, "User not found")
	}

	cartItems := make([]models.CartItem, 0)
	if err := database.Gorm.Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, cartItems)
}

func CrerateCartItem(c echo.Context) error {
	body := models.CreateCartItem{}
	c.Bind(&body)

	if err := database.Gorm.First(&models.User{}, body.UserID).Error; err != nil {
		return c.String(http.StatusNotFound, "User not found")
	}

	if err := database.Gorm.First(&models.Product{}, body.ProductID).Error; err != nil {
		return c.String(http.StatusNotFound, "Product not found")
	}

	cartItem := models.CartItem{
		UserID:    body.UserID,
		ProductID: body.ProductID,
		Quantity:  body.Quantity,
	}

	if err := database.Gorm.Create(&cartItem).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, cartItem)
}

func UpdateCartItem(c echo.Context) error {
	cartItemID := c.Param("id")

	newCartItem := models.UpdateCartItem{}
	c.Bind(&newCartItem)

	var oldCartItem models.CartItem
	if err := database.Gorm.Find(&oldCartItem, cartItemID).Error; err != nil {
		return c.String(http.StatusNotFound, "Item not found")
	}

	oldCartItem.Quantity = newCartItem.Quantity

	if err := database.Gorm.Where("id = ?", cartItemID).Save(&oldCartItem).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, oldCartItem)
}

func DeleteCartItem(c echo.Context) error {
	cartItemID := c.Param("id")

	var cartItem models.CartItem
	if err := database.Gorm.First(&cartItem, cartItemID).Error; err != nil {
		return c.String(http.StatusNotFound, "Item not found")
	}

	if err := database.Gorm.Unscoped().Delete(&cartItem, cartItemID).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "Cart Item deleted successfully")
}
