package models

import "gopkg.in/guregu/null.v4"

type Category struct {
	Model

	Name string `json:"name"`

	ImageURL  null.String `json:"image"`
	ImageName null.String `json:"-"`

	// Products	[]Product		`json:"products,omitempty"`
}

type CreateCategory struct {
	Name string `json:"name"`
}

type UpdateCategory struct {
	Name string `json:"name,omitempty"`
}
