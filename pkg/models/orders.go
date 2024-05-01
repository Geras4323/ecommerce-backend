package models

type Order struct {
	Model

	Total float64 `json:"total"`
	State uint8   `json:"state" gorm:"default:0"`

	OrderProducts []OrderProduct `json:"orderProducts"`
	Payments      []Payment      `json:"payments"`

	UserID uint `json:"userID"`
	User   User `json:"user,omitempty" gorm:"foreignKey:UserID"`

	Products uint `json:"products" gorm:"-:migration;<-:false"`
}

type CreateOrder struct {
	ProductID uint `json:"productID"`
	Quantity  uint `json:"quantity"`
}

type UpdateOrder struct {
	State uint8 `json:"state"`
}

type AddProduct struct {
	// OrderID		uint		`json:"orderID"`
	ProductID uint `json:"productID"`
	Quantity  uint `json:"quantity"`
}
