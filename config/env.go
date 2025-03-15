package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Envs = initConfig()

type Config struct {
	DB_ADDR                string
	SERVER_ADDR            string
	JWT_SECRET_KEY         string
	JWT_EXPIRATION_SECONDS int
}

func initConfig() Config {
	godotenv.Load()

	return Config{
		DB_ADDR:                getEnvStr("DB_ADDR"),
		SERVER_ADDR:            getEnvStr("SERVER_ADDR"),
		JWT_SECRET_KEY:         getEnvStr("JWT_SECRET_KEY"),
		JWT_EXPIRATION_SECONDS: getEnvInt("JWT_EXPIRATION_SECONDS"),
	}
}

func getEnvStr(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return ""
}

func getEnvInt(key string) int {
	str, ok := os.LookupEnv(key)
	if !ok {
		return 0
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return val

}
