package models

type Payment struct {
	Model

	OrderID	uint		`json:"orderID"`
	Amount	float32	`json:"amount"`
}

type CreatePayment struct {
	OrderID	uint		`json:"orderID"`
	Amount	float32	`json:"amount"`
}

type UpdatePayment struct {
	OrderID	uint		`json:"orderID"`
	Amount	float32	`json:"amount"`
}