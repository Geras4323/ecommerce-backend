package controllers

import (
	"net/http"

	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/labstack/echo/v4"
)

// GET /api/v1/suppliers //////////////////////////////////////////////////////
func GetSuppliers(c echo.Context) error {
	suppiers := make([]models.Supplier, 0)

	if err := database.Gorm.Find(&suppiers).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, suppiers)
}

// GET /api/v1/suppliers/:id
func GetSupplier(c echo.Context) error {
	supplierID := c.Param("id")

	var supplier models.Supplier
	if err := database.Gorm.First(&supplier, supplierID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, supplier)
}

// POST /api/v1/suppliers //////////////////////////////////////////////////////
func CreateSupplier(c echo.Context) error {
	body := models.CreateSupplier{}
	c.Bind(&body)

	supplier := models.Supplier{
		Name: body.Name,
	}

	if err := database.Gorm.Create(&supplier).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, supplier)
}

// PUT /api/v1/suppliers/:id //////////////////////////////////////////////////////
func UpdateSupplier(c echo.Context) error {
	supplierID := c.Param("id")

	var oldSupplier models.Supplier
	var newSupplier models.UpdateSupplier

	c.Bind(&newSupplier)

	if err := database.Gorm.First(&oldSupplier, supplierID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	oldSupplier.Name = newSupplier.Name

	if err := database.Gorm.Where("id = ?", supplierID).Save(&oldSupplier).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, oldSupplier)
}

// PATCH /api/v1/suppliers/:id
// func PatchSupplier(c echo.Context) error {

// 	return nil
// }

// DELETE /api/v1/suppliers/:id //////////////////////////////////////////////////////
func DeleteSupplier(c echo.Context) error {
	supplierID := c.Param("id")

	var supplier models.Supplier
	if err := database.Gorm.First(&supplier, supplierID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	if err := database.Gorm.Delete(&supplier, supplierID).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{"id": supplierID, "message": "Supplier deleted successfully"})
}