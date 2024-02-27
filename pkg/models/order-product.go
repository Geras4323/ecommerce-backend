package models

type OrderProduct struct {
	Model

	Quantity uint `json:"quantity"`

	OrderID uint  `json:"-"`
	Order   Order `json:"-" gorm:"foreignKey:OrderID"`

	ProductID uint    `json:"productID"`
	Product   Product `json:"-" gorm:"foreignKey:ProductID"`
}
