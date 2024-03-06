package utils

import (
	"fmt"
	"os"
)

func GetEnvVar(key string) string {
	value, ok := os.LookupEnv(key)

	if !ok {
		fmt.Printf("ENV_VAR: '%s' couln't be found", key)
	}

	return value
}
