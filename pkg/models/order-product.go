package models

type OrderProduct struct {
	Model

	Quantity uint `json:"quantity"`

	UnitID uint `json:"-"`
	Unit   Unit `json:"unit" gorm:"foreignKey:UnitID"`

	OrderID uint  `json:"-"`
	Order   Order `json:"-" gorm:"foreignKey:OrderID"`

	ProductID uint    `json:"-"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
}

type NewOrderProduct struct {
	ProductID uint `json:"productID"`
	UnitID    uint `json:"unitID"`
	Quantity  uint `json:"quantity"`
}
