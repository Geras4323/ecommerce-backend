package models

type Product struct {
	Model

	Code          string         `json:"code"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Price         float32        `json:"price"`
	OrderProducts []OrderProduct `json:"orderProducts"`

	Images []Image `json:"images"`

	CategoryID uint     `json:"categoryID"`
	Category   Category `json:"-" gorm:"foreignKey:CategoryID"`

	SupplierID uint     `json:"supplierID"`
	Supplier   Supplier `json:"-" gorm:"foreignKey:SupplierID"`
}

type CreateProduct struct {
	CategoryID  uint    `json:"categoryID"`
	SupplierID  uint    `json:"supplierID"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}

type UpdateProduct struct {
	CategoryID  uint    `json:"categoryID"`
	SupplierID  uint    `json:"supplierID"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}

type Image struct {
	Model

	Url  string `json:"url"`
	Name string `json:"-"`

	ProductID uint    `json:"-"`
	Product   Product `json:"-" gorm:"foreignKey:ProductID"`
}
