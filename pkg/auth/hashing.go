package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	// if errVer := bcrypt.CompareHashAndPassword(hash, []byte(password)); errVer != nil {
	// 	fmt.Println(errVer.Error())
	// }

	return string(hash), err
}

func VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		if err != nil {
		fmt.Println(err.Error())
	}
	return err == nil
}