package models

type State struct {
	Model

	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type UpdateState struct {
	Active bool `json:"active"`
}
