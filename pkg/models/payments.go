package models

type Payment struct {
	Model

	OrderID uint   `json:"orderID"`
	Url     string `json:"url"`
	Name    string `json:"name"`
}
