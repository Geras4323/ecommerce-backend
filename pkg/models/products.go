package models

type Product struct {
	Model

	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Position    uint    `json:"position"`
	// OrderProducts []OrderProduct `json:"orderProducts"`

	Images []Image `json:"images"`

	CategoryID uint     `json:"categoryID"`
	Category   Category `json:"-" gorm:"foreignKey:CategoryID"`

	SupplierID uint `json:"supplierID"`
}

type CreateProduct struct {
	CategoryID  uint    `json:"categoryID"`
	SupplierID  uint    `json:"supplierID"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type UpdateProduct struct {
	CategoryID  uint    `json:"categoryID"`
	SupplierID  uint    `json:"supplierID"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type UpdatePosition struct {
	ID       uint `json:"id"`
	Position uint `json:"position"`
}
