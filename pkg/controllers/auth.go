package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/geras4323/ecommerce-backend/pkg/auth"
	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/geras4323/ecommerce-backend/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GET	/api/v1/auth/session
func GetSession(baseContext echo.Context) error {
	c := baseContext.(*auth.AuthContext)
	return c.JSON(http.StatusOK, c.User)
}

// POST /api/v1/auth/login //////////////////////////////////////////////////////
func Login(c echo.Context) error {
	body := models.Login{}
	c.Bind(&body)

	var user models.User
	result := database.Gorm.Where("email = ?", body.Email).First(&user)
	err := result.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.String(http.StatusUnauthorized, err.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if isPasswordVerified := auth.VerifyPassword(user.Password, body.Password); !isPasswordVerified {
		return c.String(http.StatusUnauthorized, "Invalid credentials")
	}

	claims := &auth.JwtLoginClaims{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}

	signedToken, err := auth.SignToken(claims, utils.GetEnvVar("JWT_LOGIN_SECRET"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	cookie := new(http.Cookie)
	cookie.Name = "ec_session"
	cookie.Path = "/"
	cookie.Value = signedToken
	cookie.HttpOnly = true
	cookie.Domain = "misideaspintadas.com.ar"
	// cookie.Expires = time.Now().Add(3 * 24 * time.Hour) // expires in 3 days
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, user)
}

// POST /api/v1/auth/logout
func Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "ec_session"
	cookie.Path = "/"
	cookie.Value = ""
	cookie.MaxAge = -1
	cookie.HttpOnly = true
	cookie.Domain = "misideaspintadas.com.ar"
	c.SetCookie(cookie)

	return c.NoContent(http.StatusOK)
}

// POST /api/v1/auth/signup //////////////////////////////////////////////////////
func Signup(c echo.Context) error {
	body := models.CreateUser{}
	c.Bind(&body)

	hash, err := auth.HashPassword(body.Password)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	user := models.User{
		// Username:   body.Username,
		Email:    body.Email,
		Password: hash,
		Name:     body.Name,
		Surname:  body.Surname,
		Phone:    body.Phone,
		Role:     "customer",
	}

	if err := database.Gorm.Create(&user).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, user)
}

// POST /api/v1/auth/recovery
func RecoverPassword(c echo.Context) error {
	var body models.RecoverPassword
	c.Bind(&body)

	var user models.User
	if err := database.Gorm.Where("email = ?", body.Email).First(&user).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	claims := auth.JwtChangePasswordClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	signedToken, err := auth.SignToken(claims, utils.GetEnvVar("JWT_RES_PASS_SECRET"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	user.RecoveryToken = signedToken

	if err := database.Gorm.Where("id = ?", user.ID).Save(&user).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, signedToken)
}

// POST /api/v1/auth/change-password
func ChangePassword(c echo.Context) error {
	var body models.ChangePassword
	c.Bind(&body)

	fmt.Println(body)

	userClaims := &auth.JwtChangePasswordClaims{}
	token, err := auth.VerifyToken(body.Token, userClaims, utils.GetEnvVar("JWT_RES_PASS_SECRET"))

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var user models.User
	if err := database.Gorm.First(&user, userClaims.ID).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if token.Raw == user.RecoveryToken {
		hash, err := auth.HashPassword(body.NewPassword)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		user.Password = hash
		user.RecoveryToken = ""

		if err := database.Gorm.Where("id = ?", user.ID).Save(&user).Error; err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]any{"id": user.ID, "meesage": "Password updated successfully"})
	}

	return c.String(http.StatusUnauthorized, "Invalid password update token")
}
