package controllers

import (
	"net/http"

	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GET /api/v1/products //////////////////////////////////////////////////////
func GetProducts(c echo.Context) error {
	products := make([]models.Product, 0)
	if err := database.Gorm.Preload("Images").Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, products)
}

// GET /api/v1/products/:id
func GetProduct(c echo.Context) error {
	productID := c.Param("id")

	var product models.Product
	if err := database.Gorm.First(&product, productID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, product)
}

// POST /api/v1/products //////////////////////////////////////////////////////
func CreateProduct(c echo.Context) error {
	body := []models.CreateProduct{}
	c.Bind(&body)

	txError := database.Gorm.Transaction(func(tx *gorm.DB) error {
		for _, v := range body {
			product := models.Product{
				CategoryID:  v.CategoryID,
				SupplierID:  v.SupplierID,
				Name:        v.Name,
				Description: v.Description,
				Price:       v.Price,
			}
			if err := database.Gorm.Create(&product).Error; err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
		}
		return nil
	})

	if txError != nil {
		return txError
	}
	return c.NoContent(http.StatusOK)
}

// PUT /api/v1/products/:id //////////////////////////////////////////////////////
func UpdateProduct(c echo.Context) error {
	productID := c.Param("id")

	var newProduct models.UpdateProduct
	c.Bind(&newProduct)

	var oldProduct models.Product
	if err := database.Gorm.First(&oldProduct).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	oldProduct.Name = newProduct.Name
	oldProduct.Description = newProduct.Description
	oldProduct.Price = newProduct.Price

	if err := database.Gorm.Where("id = ?", productID).Save(&oldProduct).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, oldProduct)
}

// PATCH /api/v1/products/:id
// func PatchProduct(c echo.Context) error {

// 	return nil
// }

// DELETE /api/v1/products/:id //////////////////////////////////////////////////////
func DeleteProduct(c echo.Context) error {
	productID := c.Param("id")

	var product models.Product
	if err := database.Gorm.First(&product).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	if err := database.Gorm.Delete(&product, productID).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{"id": productID, "message": "Product deleted successfully"})
}
