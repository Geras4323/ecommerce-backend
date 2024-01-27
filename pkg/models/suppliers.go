package models

type Supplier struct {
	Model

	Name			string		`json:"name"`

	// Products	[]Product	`json:"products"`
}

type CreateSupplier struct {
	Name	string	`json:"name"`
}

type UpdateSupplier struct {
	Name	string	`json:"name"`
}