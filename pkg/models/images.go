package models

type Image struct {
	Model

	Url  string `json:"url"`
	Name string `json:"-"`

	ProductID uint    `json:"-"`
	Product   Product `json:"-" gorm:"foreignKey:ProductID"`
}
