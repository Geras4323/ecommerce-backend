package auth

import (
	"errors"
	"net/http"

	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/geras4323/ecommerce-backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AuthContext struct {
	echo.Context
	User *models.User
}

func CheckRole(r ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(baseContext echo.Context) error {
			c := baseContext.(*AuthContext)

			if validRole := utils.CheckIfInArray(r, c.User.Role); validRole {
				return next(c)
			}

			return c.String(http.StatusUnauthorized, "Unauthorized role")
		}
	}
}

func CheckAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(baseContext echo.Context) error {
		c := baseContext.(*AuthContext)

		if c.User.Role == "admin" {
			return next(c)
		}

		return c.String(http.StatusUnauthorized, "Unauthorized role - Admin needed")
	}
}

func WithAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("ec_session")
		if err != nil {
			return c.JSON(http.StatusForbidden, err.Error())
		}

		userClaims := &JwtLoginClaims{}
		_, verifyErr := VerifyToken(cookie.Value, userClaims, utils.GetEnvVar("JWT_LOGIN_SECRET"))

		if verifyErr != nil {
			return c.String(http.StatusInternalServerError, verifyErr.Error())
		}

		var user models.User
		if err := database.Gorm.Where("id = ? AND email = ?", userClaims.ID, userClaims.Email).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.String(http.StatusUnauthorized, err.Error())
			}
			return c.String(http.StatusInternalServerError, err.Error())
		}

		authContext := &AuthContext{
			Context: c,
			User:    &user,
		}

		if err := next(authContext); err != nil {
			return err
		}

		return nil
	}
}
