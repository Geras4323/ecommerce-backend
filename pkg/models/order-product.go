package models

type OrderProduct struct {
	Model

	OrderID		uint	`json:"orderID"`
	ProductID	uint	`json:"productID"`
	Quantity	uint	`json:"quantity"`
}