package models

import "gopkg.in/guregu/null.v4"

type User struct {
	Model

	// Username      string      `json:"username"`
	Email         string      `json:"email" gorm:"unique"`
	Password      string      `json:"-"`
	Name          string      `json:"name"`
	Surname       string      `json:"surname"`
	Phone         null.String `json:"phone"`
	Role          string      `json:"role" gorm:"default:customer"`
	VerifyToken   null.String `json:"-"`
	Verified      bool        `json:"verified" gorm:"default:false"`
	RecoveryToken null.String `json:"-"`

	// Orders []Order `json:"orders"`
}

type UpdateUser struct {
	// Username   string      `json:"username"`
	Name    string      `json:"name"`
	Surname string      `json:"surname"`
	Phone   null.String `json:"phone"`
	Role    string      `json:"role"`
}

type CreateUser struct {
	// Username   string      `json:"username"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Name     string      `json:"name"`
	Surname  string      `json:"surname"`
	Phone    null.String `json:"phone"`
}
