package models

import "gopkg.in/guregu/null.v4"

type Payment struct {
	Model

	OrderID  uint        `json:"orderID"`
	Url      string      `json:"url"`
	Name     null.String `json:"name"`
	Payed    null.Float  `json:"payed"`
	Received null.Float  `json:"received"`
	Platform string      `json:"platform"`
}

type MPPayment struct {
	Model

	OrderID   uint       `json:"orderID"`
	PaymentID uint       `json:"paymentID"`
	Payed     null.Float `json:"payed"`
	Received  null.Float `json:"received"`
}
