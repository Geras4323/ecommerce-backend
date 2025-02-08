package models

type CartItem struct {
	Model

	Quantity uint `json:"quantity"`

	UnitID uint `json:"-"`
	Unit   Unit `json:"unit" gorm:"foreignKey:UnitID"`

	UserID uint `json:"userID"`
	User   User `json:"-" gorm:"foreignKey:UserID"`

	ProductID uint    `json:"productID"`
	Product   Product `json:"-" gorm:"foreignKey:ProductID"`
}

type CreateCartItem struct {
	Quantity  uint   `json:"quantity"`
	Unit      string `json:"unit"`
	ProductID uint   `json:"productID"`
}

type UpdateCartItem struct {
	Quantity uint `json:"quantity"`
}
