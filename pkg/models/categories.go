package models

import "gopkg.in/guregu/null.v4"

type Category struct {
	Model

	Code string `json:"code"`
	Name string `json:"name"`

	ImageURL  null.String `json:"image"`
	ImageName null.String `json:"-"`

	// Products	[]Product		`json:"products,omitempty"`
}

type CreateCategory struct {
	Code null.String `json:"code"`
	Name string      `json:"name"`
}

type UpdateCategory struct {
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}
