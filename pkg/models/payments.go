package models

import (
	"gopkg.in/guregu/null.v4"
)

type Payment struct {
	Model

	OrderID  uint        `json:"orderID"`
	Url      null.String `json:"url" gorm:"default:null"`
	Path     null.String `json:"path" gorm:"default:null"`
	Paid     null.Float  `json:"paid" gorm:"default:null"`
	Received null.Float  `json:"received" gorm:"default:null"`
	Status   string      `json:"status" gorm:"type:ENUM('accepted','pending','rejected');not null"`
	Platform string      `json:"platform"`
}

type NewMPPayment struct {
	OrderID uint `json:"orderID"`
}

type EndMPPayment struct {
	PaymentNumber uint       `json:"paymentNumber"`
	Paid          null.Float `json:"paid"`
	Received      null.Float `json:"received"`
	Status        string     `json:"status"` // accepted | pending | rejected
}
