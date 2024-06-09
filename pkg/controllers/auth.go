package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/geras4323/ecommerce-backend/pkg/auth"
	"github.com/geras4323/ecommerce-backend/pkg/cloud"
	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/geras4323/ecommerce-backend/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
	"gopkg.in/guregu/null.v4"
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

	verificationClaims := auth.JwtVerifyEmailClaims{
		Email: body.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}

	verificationSignedToken, err := auth.SignToken(verificationClaims, utils.GetEnvVar("JWT_VERIFY_EMAIL_SECRET"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	user := models.User{
		Email:       body.Email,
		Password:    hash,
		Name:        body.Name,
		Phone:       body.Phone,
		Role:        "customer",
		VerifyToken: null.StringFrom(verificationSignedToken),
		Verified:    false,
	}

	if err := database.Gorm.Create(&user).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// SEND CONFIRMATION EMAIL
	variables := map[string]interface{}{
		"name": body.Name,
		"url":  fmt.Sprintf("%s%s%s", utils.GetEnvVar("WEB_URL"), "/sign/verifyEmail/", verificationSignedToken),
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: cloud.DefaultSender.Email,
				Name:  cloud.DefaultSender.Name,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: body.Email,
				},
			},
			Subject:          "Verificación de email",
			TemplateLanguage: true,
			TemplateID:       5934219,
			Variables:        variables,
		},
	}

	_, mailErr := cloud.SendMail(messagesInfo)
	if mailErr != nil {
		c.String(http.StatusInternalServerError, mailErr.Error())
	}

	// LOGIN USER
	loginClaims := &auth.JwtLoginClaims{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}

	loginSignedToken, err := auth.SignToken(loginClaims, utils.GetEnvVar("JWT_LOGIN_SECRET"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	cookie := new(http.Cookie)
	cookie.Name = "ec_session"
	cookie.Path = "/"
	cookie.Value = loginSignedToken
	cookie.HttpOnly = true
	cookie.Domain = "misideaspintadas.com.ar"
	// cookie.Expires = time.Now().Add(3 * 24 * time.Hour) // expires in 3 days
	c.SetCookie(cookie)

	return c.JSON(http.StatusCreated, user)
}

// POST /api/v1/signup/verify/:token
func VerifyEmail(c echo.Context) error {
	tokenParam := c.Param("token")

	claims := &auth.JwtVerifyEmailClaims{}
	token, err := auth.VerifyToken(tokenParam, claims, utils.GetEnvVar("JWT_VERIFY_EMAIL_SECRET"))

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var user models.User
	if err := database.Gorm.Where("email = ?", claims.Email).First(&user).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if token.Raw == user.VerifyToken.String {
		user.Verified = true
		user.VerifyToken = null.String{}

		if err := database.Gorm.Where("id = ?", user.ID).Save(&user).Error; err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]any{"id": user.ID, "meesage": "Account verified successfully"})
	}

	return c.String(http.StatusUnauthorized, "Invalid email verify token")
}

// POST /api/v1/auth/signup/verify/restart
func RestarEmailVerification(baseContext echo.Context) error {
	c := baseContext.(*auth.AuthContext)

	var oldUser models.User
	if err := database.Gorm.First(&oldUser, c.User.ID).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	claims := auth.JwtVerifyEmailClaims{
		Email: c.User.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}

	signedToken, err := auth.SignToken(claims, utils.GetEnvVar("JWT_VERIFY_EMAIL_SECRET"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	oldUser.VerifyToken = null.StringFrom(signedToken)

	if err := database.Gorm.Where("id = ?", c.User.ID).Save(&oldUser).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	variables := map[string]interface{}{
		"name": c.User.Name,
		"url":  fmt.Sprintf("%s%s%s", utils.GetEnvVar("WEB_URL"), "/sign/verifyEmail/", signedToken),
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: cloud.DefaultSender.Email,
				Name:  cloud.DefaultSender.Name,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: c.User.Email,
				},
			},
			Subject:          "Verificación de email",
			TemplateLanguage: true,
			TemplateID:       5934219,
			Variables:        variables,
		},
	}

	_, mailErr := cloud.SendMail(messagesInfo)
	if mailErr != nil {
		c.String(http.StatusInternalServerError, mailErr.Error())
	}

	return c.NoContent(http.StatusOK)
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

	user.RecoveryToken = null.StringFrom(signedToken)

	if err := database.Gorm.Where("id = ?", user.ID).Save(&user).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	variables := map[string]interface{}{
		"name": user.Name,
		"url":  fmt.Sprintf("%s%s%s", utils.GetEnvVar("WEB_URL"), "/sign/resetPassword/", signedToken),
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: cloud.DefaultSender.Email,
				Name:  cloud.DefaultSender.Name,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: user.Email,
				},
			},
			Subject:          "Restablecimiento de contraseña",
			TemplateLanguage: true,
			TemplateID:       6037343,
			Variables:        variables,
		},
	}

	_, mailErr := cloud.SendMail(messagesInfo)
	if mailErr != nil {
		c.String(http.StatusInternalServerError, mailErr.Error())
	}

	return c.NoContent(http.StatusOK)
}

// POST /api/v1/auth/change-password
func ChangePassword(c echo.Context) error {
	var body models.ChangePassword
	c.Bind(&body)

	claims := &auth.JwtChangePasswordClaims{}
	token, err := auth.VerifyToken(body.Token, claims, utils.GetEnvVar("JWT_RES_PASS_SECRET"))

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var user models.User
	if err := database.Gorm.First(&user, claims.ID).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if token.Raw == user.RecoveryToken.String {
		hash, err := auth.HashPassword(body.NewPassword)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		user.Password = hash
		user.RecoveryToken.String = ""

		if err := database.Gorm.Where("id = ?", user.ID).Save(&user).Error; err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]any{"id": user.ID, "message": "Password updated successfully"})
	}

	return c.String(http.StatusUnauthorized, "Invalid password update token")
}
