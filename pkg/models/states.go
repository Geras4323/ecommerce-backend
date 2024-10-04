package models

import "gopkg.in/guregu/null.v4"

type State struct {
	Model

	Name   string    `json:"name"`
	Active bool      `json:"active"`
	From   null.Time `json:"from"`
	To     null.Time `json:"to"`
}

type UpdateState struct {
	Active null.Bool `json:"active"`
	From   null.Time `json:"from"`
	To     null.Time `json:"to"`
}
