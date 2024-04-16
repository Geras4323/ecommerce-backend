package controllers

import (
	"fmt"
	"net/http"

	"github.com/geras4323/ecommerce-backend/pkg/cloud"
	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/geras4323/ecommerce-backend/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
)

// GET /api/v1/categories
func GetCategories(c echo.Context) error {
	categories := make([]models.Category, 0)

	if err := database.Gorm.Order("created_at DESC").Find(&categories).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, categories)
}

// GET /api/v1/categories/:id
func GetCategory(c echo.Context) error {
	categoryID := c.Param("id")

	var category models.Category
	if err := database.Gorm.First(&category, categoryID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, category)
}

// POST /api/v1/categories
func CreateCategory(c echo.Context) error {
	body := models.CreateCategory{}
	c.Bind(&body)

	category := models.Category{
		Code: body.Code.String,
		Name: body.Name,
	}

	if err := database.Gorm.Create(&category).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, category)
}

// POST /api/v1/categories/:id/image
func UploadCategoryImage(c echo.Context) error {
	categoryID := c.Param("id")
	var oldCategory models.Category
	if err := database.Gorm.Find(&oldCategory, categoryID).Error; err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Upload new file
	fileFolder := fmt.Sprintf("%s/%s", utils.GetEnvVar("CLOUDINARY_ENV_FOLDER"), "categories")
	fileName := uuid.New().String()
	filePath := fmt.Sprintf("%s/%s", fileFolder, fileName)

	upRes, upErr := cloud.UploadImage(file, fileFolder, fileName)
	if upErr != nil {
		fmt.Println(upErr.Error())
		return c.JSON(http.StatusBadRequest, upErr.Error())
	}

	// Replace old file
	oldFilePath := oldCategory.ImageName.String

	if oldFilePath != "" {
		_, delErr1 := cloud.DeleteImage(oldFilePath)
		if delErr1 != nil {
			_, delErr2 := cloud.DeleteImage(filePath)
			if delErr2 != nil {
				return c.JSON(http.StatusInternalServerError, delErr2.Error())
			}
			return c.JSON(http.StatusInternalServerError, delErr1.Error())
		}
	}

	oldCategory.ImageURL = null.StringFrom(upRes.SecureURL)
	oldCategory.ImageName = null.StringFrom(filePath)

	if err := database.Gorm.Where("id = ?", categoryID).Save(&oldCategory).Error; err != nil {
		return c.String(http.StatusConflict, err.Error())
	}

	return c.JSON(http.StatusOK, upRes)
}

// PUT /api/v1/categories/:id
func UpdateCategory(c echo.Context) error {
	categoryID := c.Param("id")

	var newCategory models.UpdateCategory
	c.Bind(&newCategory)

	var oldCategory models.Category
	if err := database.Gorm.First(&oldCategory, categoryID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	oldCategory.Name = newCategory.Name
	oldCategory.Code = newCategory.Code

	if err := database.Gorm.Where("id = ?", categoryID).Save(&oldCategory).Error; err != nil {
		return c.String(http.StatusConflict, err.Error())
	}

	return c.JSON(http.StatusOK, oldCategory)
}

// DELETE /api/v1/categories/:id
func DeleteCategory(c echo.Context) error {
	categoryID := c.Param("id")

	var category models.Category
	if err := database.Gorm.First(&category, categoryID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	if err := database.Gorm.Delete(&category, categoryID).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	_, err := cloud.DeleteImage(category.ImageName.String)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{"id": categoryID, "message": "Category deleted successfully"})
}

// DELETE /api/v1/categories/:id/image
func DeleteCategoryImage(c echo.Context) error {
	categoryID := c.Param("id")
	var oldCategory models.Category
	if err := database.Gorm.First(&oldCategory, categoryID).Error; err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	res, err := cloud.DeleteImage(oldCategory.ImageName.String)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	oldCategory.ImageURL = null.String{}
	oldCategory.ImageName = null.String{}

	if err := database.Gorm.Save(&oldCategory).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
