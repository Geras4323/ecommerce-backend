package models

type CartItem struct {
	Model

	Quantity uint `json:"quantity"`

	UserID uint `json:"userID"`
	User   User `json:"-" gorm:"foreignKey:UserID"`

	ProductID uint    `json:"productID"`
	Product   Product `json:"-" gorm:"foreignKey:ProductID"`
}

type CreateCartItem struct {
	Quantity  uint `json:"quantity"`
	UserID    uint `json:"userID"`
	ProductID uint `json:"productID"`
}

type UpdateCartItem struct {
	Quantity uint `json:"quantity"`
}
