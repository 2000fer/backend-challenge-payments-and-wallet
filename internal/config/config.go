package config

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var configByScope = map[string]Config{
	LocalScope: {
		ServerPort: ":8080",
		GinMode:    gin.DebugMode,
	},
	StagingScope: {
		ServerPort: ":8080",
		GinMode:    gin.TestMode,
	},
	ProductionScope: {
		ServerPort: ":8080",
		GinMode:    gin.ReleaseMode,
	},
}

type Config struct {
	ServerPort string
	GinMode    string
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
