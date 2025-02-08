package controllers

import (
	"net/http"

	"github.com/geras4323/ecommerce-backend/pkg/auth"
	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/geras4323/ecommerce-backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

var CartErrors = map[string]string{
	"Internal": "Ocurrió un error durante la carga del carrito",

	"NotFound": "Item no encontrado",
	"Create":   "Error al agregar producto al carrito",
	"Update":   "Error al actualizar carrito",
	"Delete":   "Error al eliminar item del carrito",

	"Clear": "Error al vaciar el carrito",
}

// GET /api/v1/cart //////////////////////////////////////////////////////////////////////////
func GetCartItems(baseContext echo.Context) error {
	c := baseContext.(*auth.AuthContext)

	if err := database.Gorm.First(&models.User{}, c.User.ID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(AuthErrors[utils.NotFound], err.Error()))
	}

	cartItems := make([]models.CartItem, 0)
	if err := database.Gorm.Where("user_id = ?", c.User.ID).Preload("Unit").Find(&cartItems).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(CartErrors[utils.Internal], err.Error()))
	}

	return c.JSON(http.StatusOK, cartItems)
}

// POST /api/v1/cart /////////////////////////////////////////////////////////////////////////
func CrerateCartItem(baseContext echo.Context) error {
	c := baseContext.(*auth.AuthContext)

	if err := database.Gorm.First(&models.User{}, c.User.ID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(AuthErrors[utils.NotFound], err.Error()))
	}

	body := models.CreateCartItem{}
	c.Bind(&body)

	if err := database.Gorm.First(&models.Product{}, body.ProductID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(ProductErrors[utils.NotFound], err.Error()))
	}

	var unit models.Unit
	if err := database.Gorm.Where("product_id = ? AND unit = ?", body.ProductID, body.Unit).First(&unit).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(UnitsErrors[utils.NotFound], err.Error()))
	}

	cartItem := models.CartItem{
		UserID:    c.User.ID,
		ProductID: body.ProductID,
		Quantity:  body.Quantity,
		UnitID:    unit.ID,
	}

	if err := database.Gorm.Create(&cartItem).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(CartErrors[utils.Create], err.Error()))
	}

	return c.JSON(http.StatusCreated, cartItem)
}

func UpdateCartItem(c echo.Context) error {
	cartItemID := c.Param("id")

	newCartItem := models.UpdateCartItem{}
	c.Bind(&newCartItem)

	var oldCartItem models.CartItem
	if err := database.Gorm.Find(&oldCartItem, cartItemID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(CartErrors[utils.NotFound], err.Error()))
	}

	oldCartItem.Quantity = newCartItem.Quantity

	if err := database.Gorm.Where("id = ?", cartItemID).Save(&oldCartItem).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(CartErrors[utils.Update], err.Error()))
	}

	return c.JSON(http.StatusOK, oldCartItem)
}

func DeleteCartItem(c echo.Context) error {
	cartItemID := c.Param("id")

	var cartItem models.CartItem
	if err := database.Gorm.First(&cartItem, cartItemID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(CartErrors[utils.NotFound], err.Error()))
	}

	if err := database.Gorm.Unscoped().Delete(&cartItem, cartItemID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(CartErrors[utils.Delete], err.Error()))
	}

	return c.String(http.StatusOK, "Cart Item deleted successfully")
}
