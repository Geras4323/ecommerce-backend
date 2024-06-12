package controllers

import (
	"net/http"

	"github.com/geras4323/ecommerce-backend/pkg/auth"
	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/geras4323/ecommerce-backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

var UsersErrors = map[string]string{
	"Internal": "Ocurrió un error durante la carga de los usuarios",

	"NotFound": "Usuario no encontrado",
	"Create":   "Error al crear usuario",
	"Update":   "Error al actualizar usuario",
	"Delete":   "Error al eliminar usuario",
}

// GET /api/v1/users //////////////////////////////////////////////////////
func GetUsers(c echo.Context) error {
	users := make([]models.User, 0)
	if err := database.Gorm.Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(UsersErrors[utils.Internal], err.Error()))
	}

	return c.JSON(http.StatusOK, users)
}

// GET /api/v1/users/:id
func GetUser(c echo.Context) error {
	userID := c.Param("id")

	var user models.User
	if err := database.Gorm.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(UsersErrors[utils.NotFound], err.Error()))
	}

	return c.JSON(http.StatusOK, user)
}

// PUT /api/v1/users/:id //////////////////////////////////////////////////////
func ChangeUserRole(c echo.Context) error {
	userID := c.Param("id")

	var newUser models.ChangeUserRole
	c.Bind(&newUser)

	var oldUser models.User
	if err := database.Gorm.First(&oldUser, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(UsersErrors[utils.NotFound], err.Error()))
	}

	oldUser.Role = newUser.Role

	if err := database.Gorm.Where("id = ?", userID).Save(&oldUser).Error; err != nil {
		return c.JSON(http.StatusConflict, utils.SCTMake(UsersErrors[utils.Create], err.Error()))
	}

	return c.JSON(http.StatusOK, oldUser)
}

// PATCH /api/v1/users/update-data ///////////////////////////////////////////////
func UpdateUser(baseContext echo.Context) error {
	c := baseContext.(*auth.AuthContext)

	var body models.UpdateUser
	c.Bind(&body)

	var oldUser models.User
	if err := database.Gorm.First(&oldUser, c.User.ID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(UsersErrors[utils.NotFound], err.Error()))
	}

	oldUser.Name = body.Name
	// oldUser.Email = body.Email
	oldUser.Phone = body.Phone

	if err := database.Gorm.Where("id = ?", c.User.ID).Save(&oldUser).Error; err != nil {
		return c.JSON(http.StatusConflict, utils.SCTMake(UsersErrors[utils.Update], err.Error()))
	}

	return c.JSON(http.StatusOK, oldUser)
}

// DELETE /api/v1/users/:id //////////////////////////////////////////////////////
func DeleteUser(c echo.Context) error {
	userID := c.Param("id")

	var user models.User
	if err := database.Gorm.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(UsersErrors[utils.NotFound], err.Error()))
	}

	// if err := database.Gorm.Delete(&user, userID).Error; err != nil {	// UNSCOPED ?
	if err := database.Gorm.Unscoped().Delete(&user, userID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(UsersErrors[utils.Delete], err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]any{"id": userID, "message": "User deleted successfully"})
}
