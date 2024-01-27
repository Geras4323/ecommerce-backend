package models

import "gopkg.in/guregu/null.v4"

type User struct {
	Model

	Username      string      `json:"username"`
	Email         string      `json:"email" gorm:"unique"`
	Password      string      `json:"-"`
	First_name    string      `json:"first_name"`
	Last_name     string      `json:"last_name"`
	Phone         null.String `json:"phone"`
	Role          string      `json:"role" gorm:"default:customer"`
	RecoveryToken string      `json:"-"`

	// UUID in case of needing more security
	// SecurityUUID  string      `json:"-"`

	Orders []Order `json:"orders"`
	// Active			bool				`json:"active" gorm:"default:true"`
}

type UpdateUser struct {
	Username   string      `json:"username"`
	First_name string      `json:"first_name"`
	Last_name  string      `json:"last_name"`
	Phone      null.String `json:"phone"`
	Role       string      `json:"role"`
}

type CreateUser struct {
	Username   string      `json:"username"`
	Email      string      `json:"email"`
	Password   string      `json:"password"`
	First_name string      `json:"first_name"`
	Last_name  string      `json:"last_name"`
	Phone      null.String `json:"phone"`
}
