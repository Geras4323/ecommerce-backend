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

var ProductErrors = map[string]string{
	"Internal": "Ocurrió un error durante la carga de los productos",

	"NotFound": "Producto no encontrado",
	"Create":   "Error al crear producto",
	"Update":   "Error al actualizar producto",
	"Delete":   "Error al eliminar producto",
}

var UnitsErrors = map[string]string{
	"Internal": "Ocurrió un error durante la carga de las unidades",

	"NotFound": "Unidad no encontrada",
	"Create":   "Error al crear unidad",
	"Update":   "Error al actualizar unidad",
	"Delete":   "Error al eliminar unidad",
}

// GET /api/v1/products //////////////////////////////////////////////////////
func GetProducts(c echo.Context) error {
	showAll := c.QueryParam("showAll")

	query := database.Gorm.Preload("Units").Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Order("position ASC")
	}).Order("position ASC")

	if showAll != "true" {
		query = query.Where("listed = 1")
	}

	products := make([]models.Product, 0)
	if err := query.Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(ProductErrors[utils.Internal], err.Error()))
	}

	return c.JSON(http.StatusOK, products)
}

// GET /api/v1/products/:id
func GetProduct(c echo.Context) error {
	productID := c.Param("id")

	var product models.Product
	if err := database.Gorm.Preload("Units").Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Order("position ASC")
	}).First(&product, productID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(ProductErrors[utils.Internal], err.Error()))
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
	// 			return c.JSON(http.StatusInternalServerError, utils.SCTMake(ProductErrors[utils.Create], err.Error()))
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
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(ProductErrors[utils.Create], err.Error()))
	}

	units := make([]models.Unit, 0)
	for _, u := range body.Units {
		units = append(units, models.Unit{
			Unit:      u.Unit,
			Price:     u.Price,
			ProductID: product.ID,
		})
	}

	fmt.Println(units)

	if err := database.Gorm.Create(&units).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(ProductErrors[utils.Create], err.Error()))
	}

	return c.JSON(http.StatusOK, "")
	// return c.JSON(http.StatusCreated, product)
}

// POST /api/v1/products/:id/images
func UploadProductImages(c echo.Context) error {
	productID := c.Param("id")
	var oldProduct models.Product
	if err := database.Gorm.Find(&oldProduct, productID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(ProductErrors[utils.NotFound], err.Error()))
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(utils.CommonErrors[utils.FileLoad], err.Error()))
	}
	files := form.File["images"]

	// Look for last image (for when adding images to existent product (its position))
	lastImage := models.Image{}
	database.Gorm.Select("position").Where("product_id = ?", productID).Last(&lastImage)

	// Upload and save new files to DB
	for i, file := range files {
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.SCTMake(utils.CommonErrors[utils.FileLoad], err.Error()))
		}
		defer src.Close()

		// Upload
		fileFolder := fmt.Sprintf("%s/%s", utils.GetEnvVar("CLOUDINARY_ENV_FOLDER"), "products")
		fileName := uuid.New().String()
		filePath := fmt.Sprintf("%s/%s", fileFolder, fileName)

		upRes, upErr := cloud.UploadImage(file, fileFolder, fileName)
		if upErr != nil {
			return c.JSON(http.StatusBadRequest, utils.SCTMake(utils.CommonErrors[utils.FileUpload], upErr.Error()))
		}

		// Save to DB
		var newImage models.Image

		prodID, _ := strconv.Atoi(productID)

		newImage.Url = upRes.SecureURL
		newImage.Name = filePath
		newImage.ProductID = uint(prodID)
		newImage.Position = lastImage.Position + uint(i)

		if err := database.Gorm.Create(&newImage).Error; err != nil {
			return c.JSON(http.StatusConflict, utils.SCTMake(utils.CommonErrors[utils.FileSave], err.Error()))
		}
	}

	return c.JSON(http.StatusOK, "Images uploaded succesfully")
}

// PUT /api/v1/products/:id //////////////////////////////////////////////////////
func UpdateProduct(c echo.Context) error {
	productID := c.Param("id")

	var newData models.UpdateProduct
	c.Bind(&newData)

	var productInSomeOrder bool
	database.Gorm.Raw(`
		SELECT ? IN (
			SELECT DISTINCT(product_id)
			FROM order_products op
		)
	`, productID).Scan(&productInSomeOrder)

	var oldProduct models.Product
	if err := database.Gorm.First(&oldProduct, productID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(ProductErrors[utils.NotFound], err.Error()))
	}

	// Create new version
	// if productInSomeOrder {
	// 	txError := database.Gorm.Transaction(func(tx *gorm.DB) error {
	// 		newProduct := models.Product{
	// 			CategoryID:  newData.CategoryID,
	// 			SupplierID:  newData.SupplierID,
	// 			Code:        newData.Code,
	// 			Name:        newData.Name,
	// 			Description: newData.Description,
	// 			Price:       newData.Price,
	// 			Position:    oldProduct.Position,
	// 			Listed:      true,
	// 		}

	// 		if err := database.Gorm.Create(&newProduct).Error; err != nil {
	// 			return c.JSON(http.StatusInternalServerError, utils.SCTMake(ProductErrors[utils.Create], err.Error()))
	// 		}

	// 		// Reasign images
	// 		images := make([]models.Image, 0)
	// 		if err := database.Gorm.Where("product_id = ?", productID).Find(&images).Error; err != nil {
	// 			return c.JSON(http.StatusInternalServerError, utils.SCTMake(ImageErrors[utils.NotFound], err.Error()))
	// 		}

	// 		for i, image := range images {
	// 			// Duplicate for old product
	// 			if image.Position == 0 {
	// 				newImage := models.Image{
	// 					Url:       image.Url,
	// 					ProductID: image.ProductID,
	// 					Name:      image.Name,
	// 					Position:  0,
	// 				}
	// 				if err := database.Gorm.Create(&newImage).Error; err != nil {
	// 					return c.JSON(http.StatusConflict, utils.SCTMake(utils.CommonErrors[utils.FileSave], err.Error()))
	// 				}
	// 			}
	// 			// Reasing to new product
	// 			images[i].ProductID = newProduct.ID
	// 			if err := database.Gorm.Save(&images[i]).Error; err != nil {
	// 				return c.JSON(http.StatusConflict, utils.SCTMake(utils.CommonErrors[utils.FileSave], err.Error()))
	// 			}
	// 		}

	// 		// Hide old product
	// 		oldProduct.Listed = false
	// 		if err := database.Gorm.Where("id = ?", productID).Save(&oldProduct).Error; err != nil {
	// 			return c.JSON(http.StatusInternalServerError, utils.SCTMake(ProductErrors[utils.Update], err.Error()))
	// 		}

	// 		return c.JSON(http.StatusOK, newProduct)
	// 	})

	// 	if txError != nil {
	// 		return txError
	// 	}
	// }
	if productInSomeOrder {
		newProduct := models.Product{
			CategoryID:  newData.CategoryID,
			SupplierID:  newData.SupplierID,
			Code:        newData.Code,
			Name:        newData.Name,
			Description: newData.Description,
			Price:       newData.Price,
			Position:    oldProduct.Position,
			Listed:      true,
		}

		if err := database.Gorm.Create(&newProduct).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.SCTMake(ProductErrors[utils.Create], err.Error()))
		}

		// Reasign images
		images := make([]models.Image, 0)
		if err := database.Gorm.Where("product_id = ?", productID).Find(&images).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.SCTMake(ImageErrors[utils.NotFound], err.Error()))
		}

		for i, image := range images {
			// Duplicate for old product
			if image.Position == 0 {
				newImage := models.Image{
					Url:       image.Url,
					ProductID: image.ProductID,
					Name:      image.Name,
					Position:  0,
				}
				if err := database.Gorm.Create(&newImage).Error; err != nil {
					return c.JSON(http.StatusConflict, utils.SCTMake(utils.CommonErrors[utils.FileSave], err.Error()))
				}
			}
			// Reasing to new product
			images[i].ProductID = newProduct.ID
			if err := database.Gorm.Save(&images[i]).Error; err != nil {
				return c.JSON(http.StatusConflict, utils.SCTMake(utils.CommonErrors[utils.FileSave], err.Error()))
			}
		}

		// Update cart
		cartItems := make([]models.CartItem, 0)
		if err := database.Gorm.Where("product_id = ?", productID).Find(&cartItems).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.SCTMake(ImageErrors[utils.NotFound], err.Error()))
		}

		for i := range cartItems {
			// Reasing to new product
			cartItems[i].ProductID = newProduct.ID
			if err := database.Gorm.Save(&cartItems[i]).Error; err != nil {
				return c.JSON(http.StatusConflict, utils.SCTMake(utils.CommonErrors[utils.Update], err.Error()))
			}
		}

		// Hide old product
		oldProduct.Listed = false
		if err := database.Gorm.Where("id = ?", productID).Save(&oldProduct).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.SCTMake(ProductErrors[utils.Update], err.Error()))
		}

		return c.JSON(http.StatusOK, newProduct)
	}

	// Update old version
	oldProduct.Code = newData.Code
	oldProduct.Name = newData.Name
	oldProduct.Description = newData.Description
	oldProduct.Price = newData.Price
	oldProduct.CategoryID = newData.CategoryID
	oldProduct.SupplierID = newData.SupplierID

	if err := database.Gorm.Where("id = ?", productID).Save(&oldProduct).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(ProductErrors[utils.Update], err.Error()))
	}

	return c.JSON(http.StatusOK, oldProduct)
}

// PATCH /api/v1/products/:id/images ////////////////////////////////////////////////
func UpdateProductImages(c echo.Context) error {
	productID := c.Param("id")
	var product models.Product
	if err := database.Gorm.First(&product, productID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(ProductErrors[utils.NotFound], err.Error()))
	}

	productImages := make([]models.RearrangedImage, 0)
	c.Bind(&productImages)

	// Rearrange Images
	for _, pI := range productImages {
		var oldImage models.Image
		if err := database.Gorm.First(&oldImage, pI.Id).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.SCTMake(ImageErrors[utils.NotFound], err.Error()))
		}

		if pI.IsDeleted {
			_, errDel := cloud.DeleteImage(oldImage.Name)
			if errDel != nil {
				return c.JSON(http.StatusBadGateway, utils.SCTMake(utils.CommonErrors[utils.FileDelete], errDel.Error()))
			}
			if err := database.Gorm.Unscoped().Delete(&models.Image{}, &pI.Id).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, utils.SCTMake(utils.CommonErrors[utils.FileDelete], err.Error()))
			}
		} else {
			oldImage.Position = pI.Position

			if err := database.Gorm.Where("id = ?", pI.Id).Save(&oldImage).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, utils.SCTMake(utils.CommonErrors[utils.FileSave], err.Error()))
			}
		}
	}

	return c.JSON(http.StatusOK, productImages)
}

func UpdateProductsPositions(c echo.Context) error {
	products := make([]models.UpdatePosition, 0)
	c.Bind(&products)

	// BEGIN TRANSACTION
	txError := database.Gorm.Transaction(func(tx *gorm.DB) error {
		for _, p := range products {
			var oldProduct models.Product
			if err := database.Gorm.First(&oldProduct, p.ID).Error; err != nil {
				return c.JSON(http.StatusNotFound, utils.SCTMake(fmt.Sprintf("Producto %d no encontrado", p.ID), err.Error()))
			}

			oldProduct.Position = p.Position

			if err := tx.Save(&oldProduct).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, utils.SCTMake(fmt.Sprintf("Error al guardar posición de producto %d", p.ID), err.Error()))
			}
		}

		return nil
	})

	if txError != nil {
		return txError
	}

	return c.NoContent(http.StatusOK)
}

// DELETE /api/v1/products/:id //////////////////////////////////////////////////////
func DeleteProduct(c echo.Context) error {
	productID := c.Param("id")

	var product models.Product
	if err := database.Gorm.First(&product, productID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(ProductErrors[utils.NotFound], err.Error()))
	}

	if err := database.Gorm.Delete(&product, productID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(ProductErrors[utils.Delete], err.Error()))
	}

	// Borrar las imagenes del producto borrado
	productImages := make([]models.Image, 0)
	if err := database.Gorm.Where("product_id = ?", productID).Find(&productImages).Error; err != nil {
		return c.JSON(http.StatusForbidden, utils.SCTMake(ImageErrors[utils.NotFound], err.Error()))
	}

	for _, pI := range productImages {
		_, delErr := cloud.DeleteImage(pI.Name)
		if delErr != nil {
			return c.JSON(http.StatusBadGateway, utils.SCTMake(utils.CommonErrors[utils.FileDelete], delErr.Error()))
		}
		if err := database.Gorm.Unscoped().Delete(&models.Image{}, &pI.ID).Error; err != nil {
			return c.JSON(http.StatusForbidden, utils.SCTMake(utils.CommonErrors[utils.FileDelete], err.Error()))
		}
	}

	return c.JSON(http.StatusOK, map[string]any{"id": productID, "message": "Product deleted successfully"})
}

// DELETE /api/v1/products/:id/image //////////////////////////////////////////////////////
func DeleteProductImage(c echo.Context) error {
	productID := c.Param("id")
	var product models.Product
	if err := database.Gorm.Find(&product, productID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(ProductErrors[utils.NotFound], err.Error()))
	}

	productImages := make([]models.Image, 0)
	if err := database.Gorm.Where("product_id = ?", productID).Find(&productImages).Error; err != nil {
		return c.JSON(http.StatusForbidden, utils.SCTMake(ImageErrors[utils.NotFound], err.Error()))
	}

	for _, pI := range productImages {
		_, delErr := cloud.DeleteImage(pI.Name)
		if delErr != nil {
			return c.JSON(http.StatusInternalServerError, utils.SCTMake(utils.CommonErrors[utils.FileDelete], delErr.Error()))
		}
		if err := database.Gorm.Unscoped().Delete(&models.Image{}, &pI.ID).Error; err != nil {
			return c.JSON(http.StatusForbidden, utils.SCTMake(utils.CommonErrors[utils.FileDelete], err.Error()))
		}
	}

	return c.JSON(http.StatusOK, "Images deleted successfully")
}
