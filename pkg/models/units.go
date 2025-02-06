package models

type Unit struct {
	Model

	Unit  string  `json:"unit"`
	Price float64 `json:"price"`

	ProductID uint    `json:"-"`
	Product   Product `json:"-" gorm:"foreignKey:ProductID"`
}
