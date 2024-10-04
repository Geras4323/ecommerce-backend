package controllers

import (
	"net/http"
	"time"

	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/geras4323/ecommerce-backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
)

var StateErrors = map[string]string{
	"Internal": "Ocurrió un error durante la carga de los estados",

	"NotFound": "Estado no encontrado",
	"Create":   "Error al crear estado",
	"Update":   "Error al actualizar estado",
	"Delete":   "Error al eliminar estado",
}

// POST /api/v1/states/vacation/set ////////////////////////////////////////////////////////////////
func SetVacation(c echo.Context) error {
	vacationMode := models.State{
		Name:   "vacation",
		Active: false,
	}

	if err := database.Gorm.Create(&vacationMode).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(StateErrors[utils.Create], err.Error()))
	}

	return nil
}

// GET /api/v1/states/vacation ////////////////////////////////////////////////////////////////////
func GetVacation(c echo.Context) error {
	var vacationState models.State

	if err := database.Gorm.Where("name = ?", "vacation").First(&vacationState).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(StateErrors[utils.NotFound], err.Error()))
	}

	return c.JSON(http.StatusOK, vacationState)
}

// PATCH /api/v1/states/vacation ////////////////////////////////////////////////////////////////////
func UpdateVacation(c echo.Context) error {
	body := models.UpdateState{}
	c.Bind(&body)

	oldState := models.State{}
	if err := database.Gorm.Where("name = ?", "vacation").First(&oldState).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(StateErrors[utils.NotFound], err.Error()))
	}

	if body.Active.Valid {
		oldState.Active = body.Active.Bool
		oldState.From = null.Time{}
		oldState.To = null.Time{}
	} else if body.From.Valid && body.To.Valid {
		oldState.From = body.From
		oldState.To = body.To

		now := time.Now()
		if body.From.Time.Before(now) && body.To.Time.After(now) {
			oldState.Active = true
		} else {
			oldState.Active = false
		}
	}

	// fmt.Println(oldState.Active, oldState.From, oldState.To)

	if err := database.Gorm.Where("name = ?", "vacation").Save(&oldState).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(StateErrors[utils.Update], err.Error()))
	}

	return c.NoContent(http.StatusOK)
}
