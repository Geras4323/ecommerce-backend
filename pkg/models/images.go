package models

type Image struct {
	Model

	Url      string `json:"url"`
	Name     string `json:"-"`
	Position uint   `json:"position"`

	ProductID uint    `json:"-"`
	Product   Product `json:"-" gorm:"foreignKey:ProductID"`
}

type RearrangedImage struct {
	Id        uint `json:"id"`
	Position  uint `json:"position"`
	IsDeleted bool `json:"isDeleted"`
}
