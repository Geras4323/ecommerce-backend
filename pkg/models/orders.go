package models

type Order struct {
	Model

	UserID 				uint						`json:"userID"`
	Total					float32					`json:"total"`
	Payments			[]Payment				`json:"payments"`
	OrderProducts	[]OrderProduct	`json:"orderProducts"`
}

type CreateOrder struct {
	UserID 		uint		`json:"userID"`
	Total			float32	`json:"total"`
}

type UpdateOrder struct {
	UserID 		uint		`json:"userID"`
	Total			float32	`json:"total"`
}

type AddProduct struct {
	OrderID		uint		`json:"orderID"`
	ProductID	uint		`json:"productID"`
	Quantity	uint		`json:"quantity"`
}