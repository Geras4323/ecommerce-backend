package models

import "gopkg.in/guregu/null.v4"

type User struct {
	Model

	Email         string      `json:"email" gorm:"unique"`
	Password      string      `json:"-"`
	Name          string      `json:"name"`
	Phone         null.String `json:"phone"`
	Role          string      `json:"role" gorm:"default:customer"`
	VerifyToken   null.String `json:"-"`
	Verified      bool        `json:"verified" gorm:"default:false"`
	RecoveryToken null.String `json:"-"`

	// Orders []Order `json:"orders"`
}

type UpdateUser struct {
	Name string `json:"name"`
	// Email string      `json:"email"`
	Phone null.String `json:"phone"`
}

type ChangeUserRole struct {
	Role string `json:"role"`
}

type CreateUser struct {
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Name     string      `json:"name"`
	Phone    null.String `json:"phone"`
}
