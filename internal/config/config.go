package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Mongo struct {
	URI      string
	Database string
}

type Redis struct {
	Address string
}

type Config struct {
	Mongo *Mongo
	Redis *Redis
}

var AppConfig *Config

func init() {
	LoadConfig()
}

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	AppConfig = &Config{
		Mongo: &Mongo{
			URI:      os.Getenv("MONGO_URI"),
			Database: os.Getenv("MONGO_DATABASE"),
		},
		Redis: &Redis{
			Address: os.Getenv("REDIS_ADDRESS"),
		},
	}
}
