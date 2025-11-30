package config

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type Config struct {
	ServerPort   string
	GinMode      string
	DBConnString string
}

var configByScope = map[string]Config{
	LocalScope: {
		ServerPort:   ":8080",
		GinMode:      gin.DebugMode,
		DBConnString: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
	},
	StagingScope: {
		ServerPort:   ":8080",
		GinMode:      gin.TestMode,
		DBConnString: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
	},
	ProductionScope: {
		ServerPort:   ":8080",
		GinMode:      gin.ReleaseMode,
		DBConnString: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
	},
}

const (
	LocalScope      = "local"
	StagingScope    = "staging"
	ProductionScope = "production"
	defaultScope    = "local"
)

func LoadConfig() Config {
	scope := os.Getenv("SCOPE")
	if scope == "" {
		log.Printf("Environment variable %s not set, using default: %s\n", "SCOPE", defaultScope)
		scope = defaultScope
	}

	return configByScope[scope]
}
