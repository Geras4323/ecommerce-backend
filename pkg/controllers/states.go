package controllers

import (
	"errors"
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

// Not route - Util
func SetState(state models.State) (*models.State, error) {
	if err := database.Gorm.Create(&state).Error; err != nil {
		return nil, errors.New(utils.SCTMake(StateErrors[utils.Create], err.Error()).Comment)
	}

	return &state, nil
}

// GET /api/v1/states/vacation ////////////////////////////////////////////////////////////////////
func GetVacation(c echo.Context) error {
	var vacationState models.State

	if err := database.Gorm.Where("name = ?", "vacation").First(&vacationState).Error; err != nil {
		newState, err := SetState(models.State{
			Name:   "vacation",
			Active: false,
			From:   null.Time{},
			To:     null.Time{},
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, newState)
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

// GET /api/v1/states/mercadopago ////////////////////////////////////////////////////////////////////
func GetMPPayments(c echo.Context) error {
	var mercadopagoState models.State

	if err := database.Gorm.Where("name = ?", "mercadopago").First(&mercadopagoState).Error; err != nil {
		newState, err := SetState(models.State{
			Name:   "mercadopago",
			Active: false,
			From:   null.Time{},
			To:     null.Time{},
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, newState)
	}

	return c.JSON(http.StatusOK, mercadopagoState)
}

// PATCH /api/v1/states/mercadopago ////////////////////////////////////////////////////////////////////
func UpdateMPPayments(c echo.Context) error {
	body := models.UpdateState{}
	c.Bind(&body)

	oldState := models.State{}
	if err := database.Gorm.Where("name = ?", "mercadopago").First(&oldState).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(StateErrors[utils.NotFound], err.Error()))
	}

	if body.Active.Valid {
		oldState.Active = body.Active.Bool
		oldState.From = null.Time{}
		oldState.To = null.Time{}
	}

	if err := database.Gorm.Where("name = ?", "mercadopago").Save(&oldState).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(StateErrors[utils.Update], err.Error()))
	}

	return c.NoContent(http.StatusOK)
}

// GET /api/v1/states/units ////////////////////////////////////////////////////////////////////
func GetUnits(c echo.Context) error {
	var unitsState models.State

	if err := database.Gorm.Where("name = ?", "units").First(&unitsState).Error; err != nil {
		newState, err := SetState(models.State{
			Name:   "units",
			Active: false,
			From:   null.Time{},
			To:     null.Time{},
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, newState)
	}

	return c.JSON(http.StatusOK, unitsState)
}

// PATCH /api/v1/states/units ////////////////////////////////////////////////////////////////////
func UpdateUnits(c echo.Context) error {
	body := models.UpdateState{}
	c.Bind(&body)

	oldState := models.State{}
	if err := database.Gorm.Where("name = ?", "units").First(&oldState).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(StateErrors[utils.NotFound], err.Error()))
	}

	if body.Active.Valid {
		oldState.Active = body.Active.Bool
		oldState.From = null.Time{}
		oldState.To = null.Time{}
	}

	if err := database.Gorm.Where("name = ?", "units").Save(&oldState).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(StateErrors[utils.Update], err.Error()))
	}

	return c.NoContent(http.StatusOK)
}
