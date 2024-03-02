package utils

import (
	"fmt"
	"os"
)

func GetEnvVar(key string) string {
	value := os.Getenv(key)

	if value == "" {
		fmt.Printf("ENV_VAR: '%s' couln't be found", key)
	}

	return value
}
