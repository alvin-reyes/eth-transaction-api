package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	EtherscanApiKey string
	RateLimit       int
	BurstLimit      int
	Port            string
}

// LoadConfig loads the configuration from the environment variables.
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	return &Config{
		EtherscanApiKey: getEnv("ETHERSCAN_API_KEY", ""),
		RateLimit:       getEnvAsInt("RATE_LIMIT", 1),
		BurstLimit:      getEnvAsInt("BURST_LIMIT", 5),
		Port:            getEnv("PORT", "3000"),
	}
}

// getEnv returns the value of the environment variable with the given key or the default value if the variable is not set.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt returns the value of the environment variable with the given key as an integer or the default value if the variable is not set or cannot be parsed as an integer.
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
