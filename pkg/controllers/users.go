package controllers

import (
	"net/http"

	"github.com/geras4323/ecommerce-backend/pkg/auth"
	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/labstack/echo/v4"
)

// GET /api/v1/users //////////////////////////////////////////////////////
func GetUsers(c echo.Context) error {
	users := make([]models.User, 0)
	if err := database.Gorm.Find(&users).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

// GET /api/v1/users/:id
func GetUser(c echo.Context) error {
	userID := c.Param("id")

	var user models.User
	if err := database.Gorm.First(&user, userID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

// POST /api/v1/users //////////////////////////////////////////////////////
func CreateUser(c echo.Context) error {
	body := models.CreateUser{}
	c.Bind(&body)

	hash, err := auth.HashPassword(body.Password)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	user := models.User{
		Username:   body.Username,
		Email:      body.Email,
		Password:   hash,
		First_name: body.First_name,
		Last_name:  body.Last_name,
		Phone:      body.Phone,
		Role:       "customer",
	}

	if err := database.Gorm.Create(&user).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, user)
}

// PUT /api/v1/users/:id //////////////////////////////////////////////////////
func UpdateUser(c echo.Context) error {
	userID := c.Param("id")

	var newUser models.UpdateUser
	c.Bind(&newUser)

	var oldUser models.User
	if err := database.Gorm.First(&oldUser, userID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	oldUser.Username = newUser.Username
	oldUser.First_name = newUser.First_name
	oldUser.Last_name = newUser.Last_name
	oldUser.Phone = newUser.Phone
	oldUser.Role = newUser.Role

	if err := database.Gorm.Where("id = ?", userID).Save(&oldUser).Error; err != nil {
		return c.String(http.StatusConflict, err.Error())
	}

	return c.JSON(http.StatusOK, oldUser)
}

// DELETE /api/v1/users/:id //////////////////////////////////////////////////////
func DeleteUser(c echo.Context) error {
	userID := c.Param("id")

	var user models.User
	if err := database.Gorm.First(&user, userID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	// if err := database.Gorm.Delete(&user, userID).Error; err != nil {	// UNSCOPED ?
	if err := database.Gorm.Unscoped().Delete(&user, userID).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{"id": userID, "message": "User deleted successfully"})
}
