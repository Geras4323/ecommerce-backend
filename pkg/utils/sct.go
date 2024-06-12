package utils

const (
	Internal string = "Internal"

	Unauthorized  string = "Unauthorized"
	AlreadyExists string = "AlreadyExists"
	WrongPassword string = "WrongPassword"
	PasswordHash  string = "PasswordHash"

	Email string = "Email"

	NotFound string = "NotFound"
	Create   string = "Create"
	Update   string = "Update"
	Delete   string = "Delete"

	TokenSign    string = "TokenSign"
	TokenVerify  string = "TokenVerify"
	TokenInvalid string = "TokenInvalid"

	FileLoad   string = "FileLoad"
	FileUpload string = "FileUpload"
	FileDelete string = "FileDelete"
	FileSave   string = "FileSave"
)

type ErrorDetail struct {
	Comment string `json:"comment"`
	Error   string `json:"error"`
}

var CommonErrors = map[string]string{
	Internal: "Error inesperado - Algo salió mal",

	Unauthorized: "Usuario no loggeado",

	Email: "Error al enviar email",

	FileLoad:   "Error al cargar archivo",
	FileUpload: "Error al subir archivo",
	FileDelete: "Error al eliminar archivo",
	FileSave:   "Error al guardar archivo",
}

func SCTMake(detail string, errorMessage string) ErrorDetail {
	return ErrorDetail{
		Comment: detail,
		Error:   errorMessage,
	}
}
