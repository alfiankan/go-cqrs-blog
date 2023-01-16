package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// ApplicationConfig contain all application config loaded from env or (dot)env file
type ApplicationConfig struct {
	PostgreeHost string
	PostgreeUser string
	PostgreePass string
	PostgreeDb   string
	PostgreePort int
	PostgreeSsl  string
	LogLevel     string
}

// Load load config from env
func Load(configFile ...string) ApplicationConfig {

	if err := godotenv.Load(configFile...); err != nil {
		panic(err)
	}

	postgreeDbPort, err := strconv.Atoi(os.Getenv("PG_DATABASE_PORT"))
	if err != nil {
		panic(err)
	}

	return ApplicationConfig{
		PostgreeHost: os.Getenv("PG_DATABASE_HOST"),
		PostgreeUser: os.Getenv("PG_DATABASE_USERNAME"),
		PostgreePass: os.Getenv("PG_DATABASE_PASSWORD"),
		PostgreeDb:   os.Getenv("PG_DATABASE_NAME"),
		PostgreePort: postgreeDbPort,
		PostgreeSsl:  os.Getenv("PG_DATABASE_SSL_MODE"),
		LogLevel:     os.Getenv("LOG_LEVEL"),
	}

}