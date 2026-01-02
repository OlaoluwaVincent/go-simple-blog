package env

import (
	"os"
	"strconv"
)

func GetEnvString(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}

func GetEnvInt(key string, defaultVal int) int {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultVal
	}

	value, err := strconv.Atoi(val)
	if err == nil {
		return value
	}
	return defaultVal
}
