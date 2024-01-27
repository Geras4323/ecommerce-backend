package models

type Login struct {
	Email				string			`json:"email"`
	Password		string			`json:"password"`
}

type RecoverPassword struct {
	Email				string			`json:"email"`
}

type ChangePassword struct {
	Token				string			`json:"token"`
	NewPassword	string			`json:"newPassword"`
}