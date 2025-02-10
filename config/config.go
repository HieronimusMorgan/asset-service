package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSchema   string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSchema:   os.Getenv("DB_SCHEMA"),
	}
}

func LoadRedisConfig() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv("REDIS_URL")
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Warning: .env file not found. Using system environment variables.")
	}
}

func GetJWTSecret() []byte {
	LoadEnv() // Ensure environment variables are loaded
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("❌ JWT_SECRET is missing in environment variables")
	}
	return []byte(secret) // ✅ Convert to byte slice
}
