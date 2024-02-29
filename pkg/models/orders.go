package models

type Order struct {
	Model

	UserID uint    `json:"userID"`
	Total  float64 `json:"total"`
	// Payments      []Payment      `json:"payments"`
	OrderProducts []OrderProduct `json:"orderProducts"`

	Products uint `json:"products" gorm:"-:migration;<-:false"` // llenar con SQL de DBeaver
}

type CreateOrder struct {
	ProductID uint `json:"productID"`
	Quantity  uint `json:"quantity"`
}

type UpdateOrder struct {
	// UserID uint    `json:"userID"`
	Total float64 `json:"total"`
}

type AddProduct struct {
	// OrderID		uint		`json:"orderID"`
	ProductID uint `json:"productID"`
	Quantity  uint `json:"quantity"`
}
