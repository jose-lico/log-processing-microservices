package envs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnvs() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}

func GetEnvAsBool(env string) bool {
	envValue := os.Getenv(env)

	if envValue == "" {
		return false
	}

	value, err := strconv.ParseBool(envValue)

	if err != nil {
		log.Printf("[ERROR] Error parsing env variable `%s` with value `%s` to bool: %v", env, envValue, err)
	}

	return value
}

func GetEnvAsInt(env string) int {
	envValue := os.Getenv(env)

	if envValue == "" {
		return 0
	}

	value, err := strconv.Atoi(envValue)

	if err != nil {
		log.Printf("[ERROR] Error parsing env variable `%s` with value `%s` to int: %v", env, envValue, err)
	}

	return value
}
