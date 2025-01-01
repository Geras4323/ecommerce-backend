package models

type OrderProduct struct {
	Model

	Quantity uint `json:"quantity"`

	OrderID uint  `json:"-"`
	Order   Order `json:"-" gorm:"foreignKey:OrderID"`

	ProductID uint    `json:"-"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
}

type NewOrderProduct struct {
	ProductID uint `json:"productID"`
	Quantity  uint `json:"quantity"`
}
