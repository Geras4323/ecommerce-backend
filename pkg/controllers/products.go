package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/geras4323/ecommerce-backend/pkg/cloud"
	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/geras4323/ecommerce-backend/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GET /api/v1/products //////////////////////////////////////////////////////
func GetProducts(c echo.Context) error {
	products := make([]models.Product, 0)
	if err := database.Gorm.Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Order("position ASC")
	}).Order("created_at DESC").Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, products)
}

// GET /api/v1/products/:id
func GetProduct(c echo.Context) error {
	productID := c.Param("id")

	var product models.Product
	if err := database.Gorm.Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Order("position ASC")
	}).First(&product, productID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, product)
}

// POST /api/v1/products //////////////////////////////////////////////////////
func CreateProduct(c echo.Context) error {
	// body := []models.CreateProduct{}
	// c.Bind(&body)

	// txError := database.Gorm.Transaction(func(tx *gorm.DB) error {
	// 	for _, p := range body {
	// 		product := models.Product{
	// 			CategoryID:  p.CategoryID,
	// 			SupplierID:  p.SupplierID,
	// 			Code:        p.Code,
	// 			Name:        p.Name,
	// 			Description: p.Description,
	// 			Price:       p.Price,
	// 		}
	// 		if err := database.Gorm.Create(&product).Error; err != nil {
	// 			return c.String(http.StatusInternalServerError, err.Error())
	// 		}
	// 	}
	// 	return nil
	// })

	// if txError != nil {
	// 	return txError
	// }
	// return c.NoContent(http.StatusOK)

	body := models.CreateProduct{}
	c.Bind(&body)

	product := models.Product{
		CategoryID:  body.CategoryID,
		SupplierID:  body.SupplierID,
		Code:        body.Code,
		Name:        body.Name,
		Description: body.Description,
		Price:       body.Price,
	}

	if err := database.Gorm.Create(&product).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, product)
}

// POST /api/v1/products/:id/images
func UploadProductImages(c echo.Context) error {
	productID := c.Param("id")
	var oldProduct models.Product
	if err := database.Gorm.Find(&oldProduct, productID).Error; err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	files := form.File["images"]

	// Look for last image (for when adding images to existent product (its position))
	lastImage := models.Image{}
	database.Gorm.Select("position").Where("product_id = ?", productID).Last(&lastImage)

	// Upload and save new files to DB
	for i, file := range files {
		src, err := file.Open()
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		defer src.Close()

		// // Upload
		fileFolder := fmt.Sprintf("%s/%s", utils.GetEnvVar("CLOUDINARY_ENV_FOLDER"), "products")
		fileName := uuid.New().String()
		filePath := fmt.Sprintf("%s/%s", fileFolder, fileName)

		upRes, upErr := cloud.UploadImage(file, fileFolder, fileName)
		if upErr != nil {
			fmt.Println(upErr.Error())
			return c.JSON(http.StatusBadRequest, upErr.Error())
		}

		// Save to DB
		var newImage models.Image

		prodID, _ := strconv.Atoi(productID)

		newImage.Url = upRes.SecureURL
		newImage.Name = filePath
		newImage.ProductID = uint(prodID)
		newImage.Position = lastImage.Position + uint(i)

		if err := database.Gorm.Create(&newImage).Error; err != nil {
			return c.String(http.StatusConflict, err.Error())
		}
	}

	return c.JSON(http.StatusOK, "Images uploaded succesfully")
}

// PUT /api/v1/products/:id //////////////////////////////////////////////////////
func UpdateProduct(c echo.Context) error {
	productID := c.Param("id")

	var newProduct models.UpdateProduct
	c.Bind(&newProduct)

	var oldProduct models.Product
	if err := database.Gorm.First(&oldProduct, productID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	oldProduct.Code = newProduct.Code
	oldProduct.Name = newProduct.Name
	oldProduct.Description = newProduct.Description
	oldProduct.Price = newProduct.Price
	oldProduct.CategoryID = newProduct.CategoryID
	oldProduct.SupplierID = newProduct.SupplierID

	if err := database.Gorm.Where("id = ?", productID).Save(&oldProduct).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, oldProduct)
}

func UpdateProductImages(c echo.Context) error {
	productID := c.Param("id")
	var product models.Product
	if err := database.Gorm.Find(&product, productID).Error; err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	productImages := make([]models.RearrangedImage, 0)
	c.Bind(&productImages)

	// Rearrange Images
	for _, pI := range productImages {
		var oldImage models.Image
		if err := database.Gorm.First(&oldImage, pI.Id).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		if pI.IsDeleted {
			_, errDel := cloud.DeleteImage(oldImage.Name)
			if errDel != nil {
				return c.String(http.StatusInternalServerError, errDel.Error())
			}
			if err := database.Gorm.Unscoped().Delete(&models.Image{}, &pI.Id).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
		} else {
			oldImage.Position = pI.Position

			if err := database.Gorm.Where("id = ?", pI.Id).Save(&oldImage).Error; err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
		}
	}

	return c.JSON(http.StatusOK, productImages)
}

// DELETE /api/v1/products/:id //////////////////////////////////////////////////////
func DeleteProduct(c echo.Context) error {
	productID := c.Param("id")

	var product models.Product
	if err := database.Gorm.First(&product, productID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	if err := database.Gorm.Delete(&product, productID).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Borrar las imagenes del producto borrado
	productImages := make([]models.Image, 0)
	if err := database.Gorm.Where("product_id = ?", productID).Find(&productImages).Error; err != nil {
		return c.JSON(http.StatusForbidden, err.Error())
	}

	for _, pI := range productImages {
		_, delErr := cloud.DeleteImage(pI.Name)
		if delErr != nil {
			return c.JSON(http.StatusInternalServerError, delErr.Error())
		}
		if err := database.Gorm.Unscoped().Delete(&models.Image{}, &pI.ID).Error; err != nil {
			return c.JSON(http.StatusForbidden, err.Error())
		}
	}

	return c.JSON(http.StatusOK, map[string]any{"id": productID, "message": "Product deleted successfully"})
}

func DeleteProductImage(c echo.Context) error {
	productID := c.Param("id")
	var product models.Product
	if err := database.Gorm.Find(&product, productID).Error; err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	productImages := make([]models.Image, 0)
	if err := database.Gorm.Where("product_id = ?", productID).Find(&productImages).Error; err != nil {
		return c.JSON(http.StatusForbidden, err.Error())
	}

	for _, pI := range productImages {
		_, delErr := cloud.DeleteImage(pI.Name)
		if delErr != nil {
			return c.JSON(http.StatusInternalServerError, delErr.Error())
		}
		if err := database.Gorm.Unscoped().Delete(&models.Image{}, &pI.ID).Error; err != nil {
			return c.JSON(http.StatusForbidden, err.Error())
		}
	}

	return c.JSON(http.StatusOK, "Images deleted successfully")
}
