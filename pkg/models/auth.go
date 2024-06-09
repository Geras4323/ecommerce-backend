package models

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PasswordRecovery struct {
	Email string `json:"email"`
}

type RecoverPassword struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

type ChangePassword struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
