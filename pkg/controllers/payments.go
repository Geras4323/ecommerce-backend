package controllers

import (
	"net/http"

	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/labstack/echo/v4"
)

// GET /api/v1/payments //////////////////////////////////////////////////////
func GetPayments(c echo.Context) error {
	payments := make([]models.Payment, 0)

	if err := database.Gorm.Find(&payments).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, payments)
}

// GET /api/v1/payments/:id
func GetPayment(c echo.Context) error {
	paymentID := c.Param("id")

	var payment models.Payment
	if err := database.Gorm.First(&payment, paymentID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, payment)
}

// POST /api/v1/payments //////////////////////////////////////////////////////
func CreatePayment(c echo.Context) error {
	body := models.CreatePayment{}
	c.Bind(&body)

	payment := models.Payment{
		OrderID: body.OrderID,
		Amount:  body.Amount,
	}

	if err := database.Gorm.Create(&payment).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, payment)
}

// PUT /api/v1/payments/:id //////////////////////////////////////////////////////
func UpdatePayment(c echo.Context) error {
	paymentID := c.Param("id")

	var newPayment models.UpdatePayment
	c.Bind(&newPayment)

	var oldPayment models.Payment
	if err := database.Gorm.First(&oldPayment, paymentID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	oldPayment.OrderID = newPayment.OrderID
	oldPayment.Amount = newPayment.Amount

	if err := database.Gorm.Where("id = ?", paymentID).Save(&oldPayment).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, oldPayment)
}

// PATCH /api/v1/payments/:id
// func PatchPayment(c echo.Context) error {

// 	return nil
// }

// DELETE /api/v1/payments/:id //////////////////////////////////////////////////////
func DeletePayment(c echo.Context) error {
	paymentID := c.Param("id")

	var payment models.Payment
	if err := database.Gorm.First(&payment, paymentID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	if err := database.Gorm.Delete(&payment, paymentID).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{"id": paymentID, "message": "Payment deleted successfully"})
}
